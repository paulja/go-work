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
	# cli certs
	@cp ca.pem ../cli/internal/tls
	@cp cli.pem ../cli/internal/tls
	@cp cli-key.pem ../cli/internal/tls
	# scheduler certs
	@cp ca.pem ../scheduler/internal/tls
	@cp scheduler.pem ../scheduler/internal/tls
	@cp scheduler-key.pem ../scheduler/internal/tls
	@cp worker.pem ../scheduler/internal/tls
	@cp worker-key.pem ../scheduler/internal/tls
	@cp worker-server.pem ../scheduler/internal/tls
	@cp worker-server-key.pem ../scheduler/internal/tls
	# worker certs
	@cp ca.pem ../worker/internal/tls
	@cp scheduler.pem ../worker/internal/tls
	@cp scheduler-key.pem ../worker/internal/tls
	@cp worker.pem ../worker/internal/tls
	@cp worker-key.pem ../worker/internal/tls
	@cp worker-server.pem ../worker/internal/tls
	@cp worker-server-key.pem ../worker/internal/tls
.PHONY: copy

clear:
	@rm *.{pem,csr}
.PHONY: clear
