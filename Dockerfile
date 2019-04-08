FROM golang

ARG ID_RSA

RUN mkdir /root/.ssh
RUN echo ${ID_RSA} > /root/.ssh/id_rsa
#COPY ./build/known_hosts /root/.ssh/
RUN chmod 400 /root/.ssh/id_rsa
RUN ssh-keyscan -t rsa github.com > /root/.ssh/known_hosts

### INSTALL DEP ###
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

### INSTALL GLIDE ###
RUN curl https://glide.sh/get | sh
