package internal

// refered from offical config golang/x/crypto
const (
	kexAlgoDH1SHA1                = "diffie-hellman-group1-sha1"
	kexAlgoDH14SHA1               = "diffie-hellman-group14-sha1"
	kexAlgoDH14SHA256             = "diffie-hellman-group14-sha256"
	kexAlgoECDH256                = "ecdh-sha2-nistp256"
	kexAlgoECDH384                = "ecdh-sha2-nistp384"
	kexAlgoECDH521                = "ecdh-sha2-nistp521"
	kexAlgoCurve25519SHA256LibSSH = "curve25519-sha256@libssh.org"
	kexAlgoCurve25519SHA256       = "curve25519-sha256"

	// For the following kex only the client half contains a production
	// ready implementation. The server half only consists of a minimal
	// implementation to satisfy the automated tests.
	kexAlgoDHGEXSHA1      = "diffie-hellman-group-exchange-sha1"
	kexAlgoDHGEXSHA256    = "diffie-hellman-group-exchange-sha256"
	kexAlgoDHG16EXSHA2512 = "diffie-hellman-group16-sha512"
)

var preferredKexAlgos = []string{
	kexAlgoCurve25519SHA256, kexAlgoCurve25519SHA256LibSSH,
	kexAlgoECDH256, kexAlgoECDH384, kexAlgoECDH521,
	kexAlgoDH14SHA256, kexAlgoDH14SHA1,
	kexAlgoDH1SHA1, // added manually
	kexAlgoDHGEXSHA1, kexAlgoDHGEXSHA256, kexAlgoDHG16EXSHA2512,
}

const (
	chacha20Poly1305ID = "chacha20-poly1305@openssh.com"
	aes128cbcID        = "aes128-cbc"
	tripledescbcID     = "3des-cbc"
)

// preferredCiphers specifies the default preference for ciphers.
var preferredCiphers = []string{
	"aes256-gcm@openssh.com",
	"aes128-gcm@openssh.com",
	chacha20Poly1305ID,
	"aes128-ctr", "aes192-ctr", "aes256-ctr",
	aes128cbcID, tripledescbcID, // added manually
}

// supportedMACs specifies the default preference for MACs.
var supportedMACs = []string{
	"hmac-sha2-256-etm@openssh.com",
	"hmac-sha2-256", "hmac-sha1", "hmac-sha1-96",
	"hmac-sha2-512-etm@openssh.com",
	"hmac-sha2-512",
	"umac-128-etm@openssh.com",
	"umac-128@openssh.com",
}
