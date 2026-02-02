package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
	keySize := 4096

	// Генерация приватного ключа
	fmt.Println("🔐 Генерация RSA ключа...")
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		log.Fatal("Ошибка генерации ключа:", err)
	}

	// Создание PEM блока для приватного ключа
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Запись приватного ключа в файл
	privateFile, err := os.Create("private.pem")
	if err != nil {
		log.Fatal("Ошибка создания файла:", err)
	}
	defer privateFile.Close()

	if err := pem.Encode(privateFile, privateKeyPEM); err != nil {
		log.Fatal("Ошибка записи приватного ключа:", err)
	}
	fmt.Println("✅ Приватный ключ сохранен в private.pem")

	// Вывод информации о ключе
	fmt.Printf("\n📊 Информация о ключе:\n")
	fmt.Printf("   Размер: %d бит\n", keySize)
}
