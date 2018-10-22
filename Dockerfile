FROM golang

RUN mkdir /root/.ssh
COPY ./build/id_rsa /root/.ssh/
COPY ./build/known_hosts /root/.ssh/
RUN chmod 400 /root/.ssh/id_rsa

### INSTALL DEP ###
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

### INSTALL GLIDE ###
RUN curl https://glide.sh/get | sh
