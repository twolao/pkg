package util

import (
	"io"
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
	"github.com/yeka/zip"
	"github.com/korovkin/limiter"

	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"github.com/denisbrodbeck/machineid"
)

type Backup struct {
	Src *string
	Dest *string
	Alg *string
	Pwd *string
	PwdFile *string
	Obfuscate *bool
}

var (
	// This provides obfuscation of your password only
	// Change this passphrase to have a unique encypted passwords
	// Must be 15 bytes or more
	encrypter = encrypt{passphrase: []byte{76,142,161,244,62,182,42,55,163,126,112,115,63,13,105,11,183,145,163,204,19,76,160,189,0,112,180,1,175,125}}
)

type encrypt struct {
	passphrase []byte
}

func (b *Backup) Run() {
	
	if *b.Src == "" {
		log.Fatalln("A source directory (-src) must be provided")
	}

	if *b.Dest == "" {
		log.Fatalln("A destination directory (-dest) must be provided")
	}

	if !strings.HasSuffix(*b.Src, `/`) && !strings.HasSuffix(*b.Src, `\`) {
		src := *b.Src
		src += `/`
		b.Src = &src
	}

	if !strings.HasSuffix(*b.Dest, `/`) && !strings.HasSuffix(*b.Dest, `\`) {
		dest := *b.Dest 
		dest += `/`
		b.Dest = &dest
	}

	enc := zip.AES256Encryption
	if strings.EqualFold(*b.Alg, "ZIP") {
		enc = zip.StandardEncryption
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
	}

	if *b.PwdFile != "" {
		p := filepath.Join(dir, *b.PwdFile)
		data, err := ioutil.ReadFile(p)
		if err != nil {
			log.Fatalf("Unable to read password file: " + err.Error())
		}
		if *b.Obfuscate {
			// First try to decrypt
			pwdData, err := encrypter.decrypt(data)
			if err != nil {
				// The file musn't be encrypted yet
				pwd := string(data)
				b.Pwd = &pwd
				// Encrypt the data
				encrypted, err := encrypter.encrypt(data)
				if err != nil {
					log.Fatalf("Unable to encrypt password file: " + err.Error())
				}
				err = ioutil.WriteFile(p, encrypted, 0666)
				if err != nil {
					log.Fatalf("Unable to save encrypt password file: " + err.Error())
				}
			} else {
				pwd := string(pwdData)
				b.Pwd = &pwd
			}
		} else {
			pwd := string(data)
			b.Pwd = &pwd
		}
		
	}

	type file struct {
		Path string
		Source string
		Destination string
		Size float64
	}

	limit := limiter.NewConcurrencyLimiter(5)

	start := time.Now()
	counter := 0
	totalSize := 0.0
	changesSize := 0.0
	changes := []file{}

	log.Println("Checking for changes...")
	err = filepath.Walk(*b.Src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				counter++
				relPath := path[len(*b.Src):]
				destPath := *b.Dest + relPath
				zipPath := destPath + ".zip"
				isChanged := false
				if f, err := os.Stat(zipPath); !os.IsNotExist(err) {
					if info.ModTime().After(f.ModTime()){
						isChanged = true
					}
				} else {
					isChanged = true
				}

				if isChanged {
					change := file{
						Path: relPath,
						Source: path,
						Destination: destPath,
						Size: float64(info.Size()) / 1024.0 / 1024.0,
					}
					changes = append(changes, change)
					changesSize += change.Size
				}

				totalSize += float64(info.Size()) / 1024.0 / 1024.0
				
			}
			return nil
		})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Found %v changes (%.0f MB) of %v (%.0f MB) in: %s", len(changes), changesSize, counter, totalSize, time.Now().Sub(start))
	if len(changes) == 0 {
		return
	}

	log.Println("Backing up files...")
	doneSize := 0.0
	doneCounter := 0
	for _, f := range changes {
		f := f // Intended shadowing
		limit.Execute(func() {
			percentComplete := float64(doneSize) / float64(totalSize)
			elapsed := time.Now().Sub(start).Minutes()
			total := 0.0
			if(percentComplete > 0){
				total = elapsed / percentComplete
			}
			doneCounter++
			doneSize += f.Size
			log.Printf("%0.2f%% (%0.1f of %0.1f min)    %v of %v    %s", percentComplete * 100, elapsed, total, doneCounter, len(changes), f.Path)
			err := b.ZipFile(enc, f.Source, f.Destination, *b.Pwd)		
			if err != nil {
				log.Printf("Backup failed %s: %s", f.Path, err)
			}
		})
	}

	limit.Wait()
	log.Printf("Completed %v files (%.0f MB) in: %s", len(changes), changesSize, time.Now().Sub(start))
}

func (b *Backup) ZipFile(enc zip.EncryptionMethod, src, dest, pwd string) error {

	err := os.MkdirAll(filepath.Dir(dest), 0666)
	if err != nil {
		return err
	}
	
	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	zipPath := dest + ".zip"
	fzip, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer fzip.Close()

	zipw := zip.NewWriter(fzip)
	defer zipw.Close()
	if pwd == "" {
		// If the password is blank don't encrypt
		w, err := zipw.Create(filepath.Base(dest))
		if err != nil {
			return err
		}
		_, err = io.Copy(w, fsrc)
	} else {
		w, err := zipw.Encrypt(filepath.Base(dest), pwd, enc)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, fsrc)
	}
	
	zipw.Flush()

	return nil
}


func (e *encrypt) createHash() ([]byte, error) {
	hasher := md5.New()
	id, err := machineid.ProtectedID(string(e.passphrase[4:14]))
	if err != nil {
		return nil, err
	}
	_, err = hasher.Write(append(e.passphrase, []byte(id)...))
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func (e *encrypt) encrypt(data []byte) ([]byte, error) {
	key, err := e.createHash()
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, err
}

func (e *encrypt) decrypt(data []byte) ([]byte, error) {
	key, err := e.createHash()
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("Invalid data, too small")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func (e *encrypt) encryptFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	encrypted, err := e.encrypt(data)
	if err != nil {
		return err
	}
	_, err = f.Write(encrypted)
	return err
}

func (e *encrypt) decryptFile(filename string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filename)
	return e.decrypt(data)
}
