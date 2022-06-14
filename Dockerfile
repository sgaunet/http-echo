
FROM scratch AS final
WORKDIR /
COPY http-echo .
COPY etc /etc
USER MyUser
EXPOSE 8080
CMD [ "/http-echo" ]
