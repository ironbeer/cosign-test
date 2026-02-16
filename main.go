package main

import (
	"fmt"
	"os"

	"github.com/sigstore/sigstore-go/pkg/bundle"
	"github.com/sigstore/sigstore-go/pkg/root"
	"github.com/sigstore/sigstore-go/pkg/verify"
)

func main() {
	bundlePath := os.Args[1]
	artifactPath := os.Args[2]
	certificateIdentity := os.Args[3]
	certificateIssuer := os.Args[4]

	// 1. バンドルファイルを読み込む
	bundleData, err := os.ReadFile(bundlePath)
	if err != nil {
		panic(fmt.Errorf("failed to read bundle: %w", err))
	}
	var b bundle.Bundle
	if err = b.UnmarshalJSON(bundleData); err != nil {
		panic(fmt.Errorf("failed to unmarshal bundle: %w", err))
	}

	// 2. Sigstore の公開 TUF root を取得（信頼の起点）
	trustedRoot, err := root.FetchTrustedRoot()
	if err != nil {
		panic(fmt.Errorf("failed to fetch trusted root: %w", err))
	}

	// 3. Verifier を作成
	// - WithSignedCertificateTimestamps(1): Fulcioの証明書が本物のCTログに登録されたこと
	// - WithTransparencyLog(1): Rekorエントリを1つ以上要求する
	// - WithObserverTimestamps(1):	Rekorエントリのタイムスタンプ時点では証明書が有効だったことを確認
	verifier, err := verify.NewVerifier(
		trustedRoot,
		verify.WithSignedCertificateTimestamps(1),
		verify.WithTransparencyLog(1),
		verify.WithObserverTimestamps(1),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create verifier: %w", err))
	}

	// 4. 検証対象のファイルを開く
	file, err := os.Open(artifactPath)
	if err != nil {
		panic(fmt.Errorf("failed to open plugin: %w", err))
	}
	defer file.Close()

	// 5. 検証実行
	//    - WithArtifact: 署名対象のファイル
	//    - WithCertificateIdentity: 「誰が署名したか」を指定
	_, err = verifier.Verify(&b, verify.NewPolicy(
		verify.WithArtifact(file),
		verify.WithCertificateIdentity(verify.CertificateIdentity{
			SubjectAlternativeName: verify.SubjectAlternativeNameMatcher{
				SubjectAlternativeName: certificateIdentity,
			},
			Issuer: verify.IssuerMatcher{
				Issuer: certificateIssuer,
			},
		}),
	))
	if err != nil {
		panic(fmt.Errorf("verification failed: %w", err))
	}

	fmt.Println("Verified OK")
}
