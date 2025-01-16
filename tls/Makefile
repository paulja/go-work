default:
.PHONY: default

genca:
	@cfssl gencert \
		-initca ca-csr.json | cfssljson -bare ca
.PHONY: genca

gencert:
	@cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=ca-config.json \
		-profile=server \
		scheduler-csr.json | cfssljson -bare scheduler
	@cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=ca-config.json \
		-profile=server \
		worker-server-csr.json | cfssljson -bare worker-server
	@cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=ca-config.json \
		-profile=client \
		worker-csr.json | cfssljson -bare worker
	@cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=ca-config.json \
		-profile=client \
		cli-csr.json | cfssljson -bare cli
.PHONY: gencert

copy:
	@cp *.pem ../shared/tls
.PHONY: copy

clear:
	@rm *.{pem,csr}
.PHONY: clear
