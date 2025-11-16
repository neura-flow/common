package cipher

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestRSA1(t *testing.T) {
	plain := "123456"
	publicKey := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtjG7QySNF6H8ZP0Hb25p\nN+/ZI1DWcX6FUrkOPPTa87BmeH1wEMqWSeoHcYLHJDeh7t89Wgbl332FIJvyGf2R\natqPXPMdGf0resfhgu5+Owc+I9+fwFqy+GxmAIHvrhvcgF9/CY9Goa5CCTgdl9DH\ng0iVo/dGO7nT8sQZ2ItBwVTE7xXJcsXB77p0d7dkxjAIdQ68ZovMUm8KYbHzU/17\nZMScLjBnK+wv8zOA/rCvYXzbpCU6zuqdzvgIvA86b2MWlby1M65Ax1xmyBl7Oarn\n45VPPenVbQ9CWR8oA3QvoZngyR+u8hg4+TlBDbJyRqT7AvzTYLRwzOV2YkskpZFk\nHwIDAQAB\n-----END PUBLIC KEY-----\n"
	privateKey := "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAtjG7QySNF6H8ZP0Hb25pN+/ZI1DWcX6FUrkOPPTa87BmeH1w\nEMqWSeoHcYLHJDeh7t89Wgbl332FIJvyGf2RatqPXPMdGf0resfhgu5+Owc+I9+f\nwFqy+GxmAIHvrhvcgF9/CY9Goa5CCTgdl9DHg0iVo/dGO7nT8sQZ2ItBwVTE7xXJ\ncsXB77p0d7dkxjAIdQ68ZovMUm8KYbHzU/17ZMScLjBnK+wv8zOA/rCvYXzbpCU6\nzuqdzvgIvA86b2MWlby1M65Ax1xmyBl7Oarn45VPPenVbQ9CWR8oA3QvoZngyR+u\n8hg4+TlBDbJyRqT7AvzTYLRwzOV2YkskpZFkHwIDAQABAoIBADOKphs8f92rqac4\nHZ5ccc+tPpRLDh9VV4orZ+e+rSI7SQHVGprldNa8KhbmlEeepWTaKpUJVoZ/D+ZF\nt5u6rCS6Z8w3yofLoz08xoMvzO4OAnpLjPnxrqeworqKB7ANmbeHTHz711Nt5KiP\nA3ArVAXDxvF3xpqm21rWNymXW6bBRi38yiCE9G7ja6ZDkP9W20tvSM66k61yNeZH\nUF9tryypmw3F7OEz+PMOnvPKqcCyAzb7bQTsHGXyf2kxyZj6bZEkzUHoilwHGnj/\nepz0ddV/LKWB8jsz3GN1k+JX0J/doOyyuoNWp7dLuCbBm00iMNOYbfcNvIyQt8hB\npUeEsaECgYEA3BdZU5vNlaAHqm273+n844ZElaU84KXs3CPHRZjt46JsyAKMDRJM\n9yXYrnTKrrMuSlIfJ9G2BXLeMrEvxJyWghzCOgELpf4n/BePPBE+SwcfOD6KCjG4\nDYQwf4wz+fiXdr/Lz3BuwN85C4e7Cf3G7jdk/3MspPq1Lro9gvLT3k0CgYEA0+uG\nAAMOo6twzqe8WAblZYXCFYzYsOgiB+d2RLJqlD9EVfleD3BE41EInAvyhNwY8E47\nQzOmBfd7zH3ZZsMc6EPsTwex9JXRoKZZtAXOgB0Rl+ZDWmCJv84AST8/8pzwEPhu\nc1V11yD3xCc4Kn9cVp2fbaU8aPq2GiHDw+0JuhsCgYBKRKrE9udZ3UWY8iyas5e9\no1pTcQ3o9LTH2F7vElr8HJw+pfVil9FW+PN0cz7N0vME609OHYshrZBjZL0syHZV\nc6Tq891dZzVQ8RZJe7wcj0uurBPiusJT9U50S/hiGsvpq3D4EAWfmfPi+ytXhMZz\nLkgrl07yYRNwsDH/lTd/ZQKBgQCHHMq6hzh2MYAiwd7bYMoxCC7N/pbJc7b+wxws\ngHRjQFMZXXwS68mABNIwa42cF5fu3nH6Tpuzgi50GmjZk9yCWYv4dzeGcV7NxkG6\n/VjDZcUpy611mcc5euXDzYe/7z9AEqSY9AvFtUdC0J6GudztfGGBTrBNXktsLcra\nx+5DsQKBgQCqlkLIznmE8Uqgr/z+GJymAQ5p1oRRqgs5Se9BznVDQRra29N6UW8T\nXGs7Jq5HFx8AXXfF0qu2buF0YSQXN1UVgMPRC4tNqk9M+fMsTcNjDyNNHrEMEiq/\nffueILceub2OdcEwM9GzEudo/nrAAOk69X+7fg+rE08gLMcaSyGMeA==\n-----END RSA PRIVATE KEY-----\n"
	rsa := NewRSA(&RSAConfig{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Format:     FormatPkcs1,
	})
	cipher, err := rsa.Encrypt([]byte(plain))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", string(cipher))

	s := base64.StdEncoding.EncodeToString(cipher)
	s = "SnpDNk9nR0swTHJCUmJLNGJFcm03K21WN2pnVXV1Q1AwbXZpbFdqSU5XbEFKcCtJZVMxS0laVmxITHNzY0R2eWFVM2UrdFVkcW82THB0SGhTRVp4STl0WXZDTEYzZkw4dHYyV1pBVEdYUWZLUzFqV1h1YWlrYmZYVzk4R1VHaWFDMmpQcVg0ZGtaK1J1T1RuYzZxRS93SnNxT2Y4NnZ6bkpOcUIwUXowNEZYV1MwT01OZEVEdU0vWkdEc3pQY2FGVnNKMVBrTVU3ODZML2tCTGFYWkVUQy9Bb0dNODcwY1pjTnhwVGRYRzhvSktkSmx0WFRhM2pWUDQ3d1hZKytteDFRZGhweHRoK1EwQTcvUm5CM1h5MGxUd1doTFREUGVhblBtaFRhdFhCQ3NNOUpacEtqR1hWWG1jVTRvc1pGTzZYSEZ1eVR6M2t4clRUUmZoRC9nUVdBPT0="
	fmt.Printf("base64 enc: %s \n", s)

	ds, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	plainBytes, err := rsa.Decrypt(ds)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("dec %s\n", string(plainBytes))
}

func TestRSA2(t *testing.T) {
	plain := "hello1234"
	publicKey := "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCqfWLHwNKjRBqW2zBeqZId2/8b\nyjHfOt6an4HHV/hFLSy9x1xOJS3tFC+aqOmRFQTBIbwWcrw626GoMRQiEZ82wHEw\n1YomyJUV+lINJHxxdFgpsAlhCWY8i3coZFADo6g2s3qT8sa/NYdjJmXj/tITSrXN\nl24yRciAi8OmrZInrwIDAQAB\n-----END PUBLIC KEY-----"
	privateKey := "-----BEGIN PRIVATE KEY-----\nMIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAKp9YsfA0qNEGpbb\nMF6pkh3b/xvKMd863pqfgcdX+EUtLL3HXE4lLe0UL5qo6ZEVBMEhvBZyvDrboagx\nFCIRnzbAcTDViibIlRX6Ug0kfHF0WCmwCWEJZjyLdyhkUAOjqDazepPyxr81h2Mm\nZeP+0hNKtc2XbjJFyICLw6atkievAgMBAAECgYEAiggiJZ/T6iOFc4Xaz4lyp4Iq\nNRGq9xTujTl/FSn/8+HzS2NRNGOdn0iskgcXs0LVKphpc7NI+k4/v7CcoEisIZFy\nTZ1ZUItBxmvLaSRBvKJJ9N1ZuLRpfAhDkERGubBVj++kzWQUr7ppNGlx83wRtFHJ\nhvBjCfBloFT79Iiil0ECQQDIgebXJhJUm/ZWc00yvB1R6IbwBaLeeTrza6vENtg9\nbLLDv1qHRIXt2dfgrfghm46GY4GYLcHjWGvoq3TtHR35AkEA2ay2W9dmqeISliPH\nB1U3Jb40pzP67f0Uf5hY06GacyWl/kKGIPxPPxrNAsw4JJgIvDnwKeSLsfiARhd5\nE+P85wJBALB4JAMXrupomdZchIUyq1t7m8eELmQ/rnKvQO3gl1D4ah1+PN7woC9G\nm4lTlB+AGWCOE3EsVIkTOWX+AVrvVYECQQCVhCTegN5r4nWR25FiYA45RqU0FGhQ\nAH6MBkE9XMuSPFIAjAFFtwlX9zjKqywFNskJQWLN48ZwwJibjJQGLZwRAkAxDOAv\nJ/ab3gQ+vQxtb1i0M8vz1EUgW3o7NOfih2Z3M437jZhh4SOanzg47CK/SDdBmaE3\nioY1UCDglVObHsJB\n-----END PRIVATE KEY-----"
	rsa := NewRSA(&RSAConfig{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Format:     FormatPkcs1,
	})
	cipher, err := rsa.EncryptToBase64(plain)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("base64 enc: %s \n", cipher)

	plainStr, err := rsa.DecryptBase64(cipher)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("dec %s\n", plainStr)
}

func TestRSAEnc(t *testing.T) {
	plain := "hello1234"
	publicKey := "-----BEGIN PUBLIC KEY-----\nMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAKoR8mX0rGKLqzcWmOzbfj64K8ZIgOdH\nnzkXSOVOZbFu/TJhZ7rFAN+eaGkl3C4buccQd/EjEsj9ir7ijT7h96MCAwEAAQ==\n-----END PUBLIC KEY-----"
	rsa := NewRSA(&RSAConfig{
		PublicKey: publicKey,
		Format:    FormatPkcs1,
	})
	cipher, err := rsa.EncryptToBase64(plain)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("base64 enc: %s \n", cipher)
}

func TestRSADec(t *testing.T) {
	cipher := "ZZFhcg3BK1mf7UJZwQiDPHSv94ivTYVTC79vQjfR1OdtLyb8gW+HNG2JCW6g/XdOkn+tE9VFUsc0+7TfYOutsMnmikjqkp8NA0Uf5imucV2Z56hMBKPees3Pbt9teRIOFXT74D8k7h1wL4D0OW7vYH5wgIIWK4FsX4s/1/jB/PE="
	privateKey := "-----BEGIN PRIVATE KEY-----\nMIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEAqhHyZfSsYourNxaY\n7Nt+PrgrxkiA50efORdI5U5lsW79MmFnusUA355oaSXcLhu5xxB38SMSyP2KvuKN\nPuH3owIDAQABAkAfoiLyL+Z4lf4Myxk6xUDgLaWGximj20CUf+5BKKnlrK+Ed8gA\nkM0HqoTt2UZwA5E2MzS4EI2gjfQhz5X28uqxAiEA3wNFxfrCZlSZHb0gn2zDpWow\ncSxQAgiCstxGUoOqlW8CIQDDOerGKH5OmCJ4Z21v+F25WaHYPxCFMvwxpcw99Ecv\nDQIgIdhDTIqD2jfYjPTY8Jj3EDGPbH2HHuffvflECt3Ek60CIQCFRlCkHpi7hthh\nYhovyloRYsM+IS9h/0BzlEAuO0ktMQIgSPT3aFAgJYwKpqRYKlLDVcflZFCKY7u3\nUP8iWi1Qw0Y=\n-----END PRIVATE KEY-----"
	rsa := NewRSA(&RSAConfig{
		PrivateKey: privateKey,
		Format:     FormatPkcs1,
	})
	plainStr, err := rsa.DecryptBase64(cipher)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("dec %s\n", plainStr)
}

func TestGenRSAKey(t *testing.T) {
	pub, pri, err := GenRSA(RSABits2048, false)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("\"%s\"", strings.Replace(string(pub), "\n", "\\n", -1))
	fmt.Printf("\n")
	fmt.Printf("\"%s\"", strings.Replace(string(pri), "\n", "\\n", -1))
}
