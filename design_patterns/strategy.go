package main

import "fmt"

type PasswordProtector struct {
	user          string
	passwordName  string
	hashAlgorithm HashAlgorithm
}

type HashAlgorithm interface {
	Hash(p *PasswordProtector)
}

func NewPaswordProtector(user string, passwordName string, hashAlgorithm HashAlgorithm) *PasswordProtector {
	return &PasswordProtector{
		user:          user,
		passwordName:  passwordName,
		hashAlgorithm: hashAlgorithm,
	}
}
func (p *PasswordProtector) SetHashAlgorithm(hash HashAlgorithm) {
	p.hashAlgorithm = hash
}

func (p *PasswordProtector) Hash() {
	p.hashAlgorithm.Hash(p)
}

type SHA struct {
}

func (SHA) Hash(p *PasswordProtector) {
	fmt.Printf("Hashing using SHA %s\n", p.passwordName)
}

type MD5 struct {
}

func (MD5) Hash(p *PasswordProtector) {
	fmt.Printf("Hashing using MD5 %s\n", p.passwordName)
}

func main() {
	sha := &SHA{}
	md5 := &MD5{}

	PasswordProtector := NewPaswordProtector("angelo", "passwordgmail", sha)
	PasswordProtector.Hash()
	PasswordProtector.SetHashAlgorithm(md5)
	PasswordProtector.Hash()
}
