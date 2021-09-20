package main

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//var ZipDir string
//var EncZip string
var ZipDir string
var EncZip string

var key []byte

func main() {

args:=os.Args
l := 1

for l < len(args) {
	ZipDir += args[l] + " "
	l++
}
ZipDir = ZipDir[:len(ZipDir)-1]

//keyin := []byte(args[2])
keyin := ""
fmt.Print("PASSWORD: ")
fmt.Scanln(&keyin)

keytmp := make([]byte, 32)
copy(keytmp, keyin)

hash := sha512.New()

for i:=0;i<32;i++ {
	hash.Write(keytmp)
	keytmp = hash.Sum(keytmp)
	hash.Write(keytmp)
	keytmp = hash.Sum(keytmp)
}

key = keytmp[0x30:0x50]
//fmt.Print(hex.Dump(key))
EncZip = ZipDir

info,e := os.Stat(ZipDir)
if e!=nil {
	fmt.Println(e)
}

if info.IsDir() {
	EncZip = ZipDir + ".zl"
	Enc()
} else {
		name := strings.Index(EncZip, ".zl")
		if name == -1 {
			fmt.Println("IT IS NOT ENCRYPTED FILE! File format is not recognized. ")
			stop:=""
			fmt.Scan(&stop)
			return
		}
		Dec()
	}

}

func Enc(){

	fmt.Println("Starting enc...")
	//args:=os.Args
	//testDirectory = args[1] - WHEN I ADDED TO THE REGEDIt
	//fmt.Println(testDirectory)
	compress(ZipDir,EncZip) // 
	//decompress(testDirectory, zipDirectory)
	//decompress(testDirectory, zipDirectory)
	fileSize := ReadFile(EncZip) // % 4096
	//fmt.Println(fileSize)
	Encrypt(fileSize, EncZip) //
	fmt.Println("Encrypted!")
}
func Dec() {

	fmt.Println("Starting Dec...")
	fileSize := ReadFile(EncZip)
	Decrypt(fileSize, EncZip)
	//ZipDir = strings.Replace(ZipDir,".zl","",1)
	//decompress(ZipDir, EncZip)
	slash := strings.LastIndex(EncZip,"\\")
	ZipDir = EncZip[:slash]
	decompress(ZipDir, EncZip)
	err := os.Remove(EncZip)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Decrypted!")
}
func compress (testDirectory string, zipDirectory string){

	zipfile, err := os.OpenFile(zipDirectory, os.O_RDWR|os.O_CREATE, 0666)
	//zipfile, err := os.Create(zipDirectory)
	if err !=nil{
		fmt.Println("\n")
		fmt.Println("USE OTHER PLACE! YOU ARE NOT ADMIN!")
		fmt.Println("\n OR")
		fmt.Println("\n RUN CMD AS ADMIN: usage ---- cmd> ZeroCrypt.exe C:\\CryptedDir")
		log.Fatal(err)
		return

	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	info, err := os.Stat(testDirectory)
	if err != nil {log.Fatal(err)}
	var baseDir = filepath.Base(testDirectory)
	if info.IsDir() {
		baseDir = filepath.Base(testDirectory)
	}
	filepath.Walk(testDirectory,func(path string,info os.FileInfo, err error) error {
		if err != nil {
			return err

		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if baseDir !=""{
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, testDirectory))
		}
		if info.IsDir() {
			header.Name +="/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir(){
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_,err = io.Copy (writer, file)
		return err
	})
}
func decompress (testDirectory string, zipDirectory string){
	reader, err := zip.OpenReader(zipDirectory)
	if err != nil {
		fmt.Println("\n")
		fmt.Println("USE OTHER PLACE! YOU ARE NOT ADMIN!")
		fmt.Println("\n OR")
		fmt.Println("\n RUN CMD AS ADMIN: usage ---- cmd> ZeroCrypt.exe C:\\CryptedDir")
		log.Fatal(err)
		return
	}

	err = os.MkdirAll(testDirectory, 0777)
	if err != nil {
		fmt.Println("\n")
		fmt.Println("USE OTHER PLACE! YOU ARE NOT ADMIN!")
		fmt.Println("\n OR")
		fmt.Println("\n RUN CMD AS ADMIN: usage ---- cmd> ZeroCrypt.exe C:\\CryptedDir")
		log.Fatal(err)
		log.Fatal(err)
		return
		}
	for _, file := range reader.File {
		path := filepath.Join(testDirectory, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}
		fileReader, err := file.Open()
		if err != nil {fmt.Println(err)}
		defer fileReader.Close()

		targetFile, err:= os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err !=nil{fmt.Println(err)}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil{fmt.Println(err)}

	}
	reader.Close()
}
func ReadFile(zipDirectoryOut string) int64{

	file, err := os.Open(zipDirectoryOut)
	if err != nil {
		fmt.Println("\n")
		fmt.Println("USE OTHER PLACE! YOU ARE NOT ADMIN!")
		log.Fatal(err)
		fmt.Println(err)
		return 0
	}

	file_info, err := file.Stat()
	if err != nil{
		fmt.Println(err)
	}
	file.Close()
	lenght := file_info.Size()
	padding := 0
	for lenght % 4096 != 0 {lenght++; padding++}
	var rnd int64 = 1
	rndByte := make ([]byte, 1)
	rand.Seed(lenght)
	rand.Read(rndByte)
	i := int(rndByte[0])
	if err != nil {
		fmt.Println(err)}

	for i >= 0 {
		time := time.Now()
		t := time.UnixNano()
		rand.Seed(rnd*t)
		rnd = int64(rand.Uint64())
		i--
		//fmt.Println(rnd)
	}

	rndstate := uint64(rnd)

	if padding != 0 {
		randBuff := make([]byte, padding)
		rand.Seed(int64(rndstate))
		rand.Read(randBuff)
		file, err = os.OpenFile(zipDirectoryOut, os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			fmt.Println(err)
		}
		file.Write(randBuff)
		file.Close()
	}

	return lenght
}
func Encrypt(lenght int64, zipDirectoryOut string) {

	var shift int64
	shift = 0

	if lenght >= 4096 {

		//randBuff2 := make([]byte, 4096)
		//randBuff3 := make([]byte, 4098)
		text := make([]byte, 4096)
		//lenghtforRND := lenght
		//lenghtforShift := lenght

		file, err := os.OpenFile(zipDirectoryOut, os.O_RDWR, 0777)
		if err != nil {
			fmt.Println(err)
		}
		//key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574") // HERE YOU CAN CHANGE THE SECRET FOR YOUR PURPOSE!
		iv := []byte("MicrosoftWindows")
		block, err := aes.NewCipher(key)
		ciphertext := make([]byte, 4096)

		for lenght > -1 {

			file.ReadAt(text, shift)
			mode := cipher.NewCBCEncrypter(block, iv) //ENCRYPT
			mode.CryptBlocks(ciphertext, text)
			file.WriteAt(ciphertext, shift)
			lenght = lenght - 4096
			shift = shift + 4096
			//rand.Seed(int64(rndstate) + lenghtforRND + 1)
			//rand.Read(randBuff2)
			//lenghtforRND = lenghtforRND - 4096 //
			//file.WriteAt(randBuff2, lenghtforShift)
			//lenghtforShift = lenghtforShift + 4096
		}

		file.Close()
		//fmt.Println("11")
	}
}
func Decrypt(lenght int64, zipDirectoryOut string) {

	var shift int64
	shift = 0
		//key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574") // HERE YOU CAN CHANGE THE SECRET FOR YOUR PURPOSE!
		iv := []byte("MicrosoftWindows")
		block, _ := aes.NewCipher(key)
		ciphertext := make([]byte, 4096)
		//shift := 0
	file, err := os.OpenFile(zipDirectoryOut, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
	}
	text := make([]byte, 4096)
		for lenght > -1 {

			file.ReadAt(text, shift)
			mode := cipher.NewCBCDecrypter(block, iv) //DECRYPT
			mode.CryptBlocks(ciphertext, text)
			file.WriteAt(ciphertext, shift)
			lenght = lenght - 4096
			shift = shift + 4096
			//rand.Seed(int64(rndstate) + lenghtforRND + 1)
			//rand.Read(randBuff2)
			//lenghtforRND = lenghtforRND - 4096 //
			//file.WriteAt(randBuff2, lenghtforShift)
			//lenghtforShift = lenghtforShift + 4096
		}

		file.Close()
		//fmt.Println("11")
}
