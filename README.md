# Overview

Create encrypted files, decrypt them securely in Kubernetes/Docker. An encryption and decryption tool is provided by this package.

# How it works

You create a file that maps environment variable names to values.

# Tutorial

 1. Create a secrets file, such as secrets.plaintext:
    ```
    DB_PASSWORD=passw0rd
    AWS_ACCESS_ID=XXXXXX
    AWS_SECRET=XXXXX
    ```
 1. Encrypt your secrets using:
    
    ```
    envcrypt encrypt aes256 -keyBase64 SECRET_BASE64_ENCODED -in secrets.plaintext -out secrets.enc`
    ```
    
    The secret must be exactly 32 bytes. You can generate one using ruby with: `ruby -e 'require "securerandom"; puts SecureRandom.base64(32)'`
 1. Store them in your image `COPY secrets.enc secrets.enc`
 1. Tell [Kubernetes](https://kubernetes.io/docs/concepts/configuration/secret/) about your Key `kubectl create secret generic decryption-key --from-literal=key=SECRET`
 1. Have Kubernetes [inject the key as an environment variable](https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets-as-environment-variables) in your deployment: 
    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: secret-env-pod
    spec:
      containers:
      - name: mycontainer
        image: myimage
        env:
          - name: KEY_BASE64
            valueFrom:
              secretKeyRef:
                name: decryption-key
                key: key 
    ```
 1. execute the utility in your image and have the utility execute your program
    ```dockerfile
    FROM alpine:latest
    COPY myserver .
    COPY secrets.enc .
    ENTRYPOINT ["./envcrypt", "decrypt-to-env", "-in=secrets.enc", "-exec=./myserver", "-execArg=someargument", "-execArg=someotherargument"]
    ```

When **envcrypt** runs, it will read the files specified by in (in order you specify them), decrypt them, and set the values contained therein as environment variables. The last matching environment variable overwrites the previous ones. Any existing environment variables are preserved except for the `KEY_BASE64` environment variable and any environment variables set in your secrets.plaintext file. The `KEY_BASE64` is unset before your application is run, thereby preserving the Key encrypting key.

When the above docker image is run, the image will contain the at least following environment variables in addition to any already existing environment variables:

 * DB_PASSWORD=passw0rd
 * AWS_ACCESS_ID=XXXXXX
 * AWS_SECRET=XXXXX
 
Note, the KEY_BASE64 variable will not exist as it will be cleared when myserver is executed.

Each line in the file will be extracted as an environment variable.

# Building

 1. [Install go](https://golang.org/doc/install)
 1. `cd cmd/cli`
 1. `go build -o envcrypt -i .`

# Future work

 * Asymmetric keys
 * More encryption ciphers
 * Integration with AWS KMS and Google KMS
 * Allow the KEY_BASE64 environment variable to be changed?

# Getting help

The commandline client includes documentation. Just call `./envcrypt -h` to get more information, e.g.:

```
./envcrypt encrypt aes256 -h
```

# What the?

How does the environment variable injection work? Exec. Go will exec the process you specify whereby it replaces its own running process with your code

# Purpose

With images in Kubernetes, we need a way to inject credentials into images without using Consul/Vault or other run-time dependencies. There are many alternatives. Google opts for baking secrets into images using Cloud Build. Docker images are built with secrets and K8s can pull the images using the access credential. While I do like this idea, it means that any deployment run in Kubernetes that can read that repository can find a way to extract the secrets. I'm not aware of any official stance from Amazon or Azure, but they may exist (in secret).

We also wanted a way to rotate the credentials and have them roll themselves out. K8s deployments seemed to be a good lead. If the credentials were encrypted and part of the build process, a change to the encrypted contents would trigger a rollout deployment. If something was wrong, a rollback could also be done and the system could be restored.

There is a concern regarding combining code with configuration. If this is a concern, this utility can be used if you keep your configuration somewhere the build system can inject it into the image, still encrypted. This should cover most of the bases. However, I'm of the opinion that your configuration should also be versioned. This is a way to enforce that.

# Copyright

Chris Wojno 2019. All rights Reserved

# License

![cc logo](https://i.creativecommons.org/l/by-sa/4.0/88x31.png)

This work is licensed under a Creative Commons Attribution-ShareAlike 4.0 International License.

# Disclaimer

This is not production code. Don't use it in production. Use at your own risk. Author and contributors are not responsible for anything. You've been warned.
