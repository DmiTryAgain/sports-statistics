FROM mysql:5.7

ARG USER_ID='1000'
ARG GROUP_ID='1000'

ADD my.cnf /etc/mysql/conf.d/my.cnf

RUN groupmod -g $USER_ID mysql
RUN usermod -u $USER_ID mysql

CMD ["--character-set-server=utf8mb4", "--collation-server=utf8mb4_general_ci", "--skip-character-set-client-handshake"]

EXPOSE 3306
