import socket
import ssl
import pprint
import time
import sys

def foo(self):
    if 3 == sys.version_info[0]:
        def sendStr(msg):
            return self.send(bytes(msg, encoding='utf-8'))
        return sendStr
    else:
        return self.send

def printB(buf):
    if 3 == sys.version_info[0]:
        print(str(buf, encoding='utf-8'))
    else:
        print(buf)

def sslClient(server_name, port):
    context = ssl.create_default_context()
    sockfd = context.wrap_socket(socket.socket(socket.AF_INET), server_hostname=server_name)
    sockfd.connect((server_name, port))
    return sockfd

if "__main__" == __name__:
    BUFSIZ = 4096
    sockfd = sslClient("imap-mail.outlook.com", 993)
    # sockfd = sslClient("imap-mail.outlook.com", 993)
    sockfd.sendStr = foo(sockfd)
    print("connected!")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    sockfd.sendStr("1 CAPABILITY\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    exit(0)

    #sockfd.sendStr("2 LOGOUT\r\n");
    #buf = sockfd.recv(BUFSIZ)
    #printB(buf)
    sockfd.sendStr("2 LOGIN james@watsonserve.com 123456\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    sockfd.sendStr("3 CAPABILITY\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    sockfd.sendStr("4 SELECT INBOX\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    print("\r\n")
    sockfd.sendStr("5 UID SEARCH 1:* NOT DELETED\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    #sockfd.sendStr("6 UID FETCH 126:127 (INTERNALDATE UID RFC822.SIZE FLAGS BODY.PEEK[HEADER.FIELDS (date subject from content-type to cc bcc message-id in-reply-to references)])\r\n")
    #buf = sockfd.recv(BUFSIZ)
    #printB(buf)
    sockfd.sendStr("6 NOOP\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    sockfd.sendStr("7 UID FETCH 126 (BODYSTRUCTURE BODY.PEEK[HEADER])\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    sockfd.sendStr("8 UID FETCH 1:3 (INTERNALDATE UID RFC822.SIZE FLAGS BODY.PEEK[HEADER.FIELDS (date subject from content-type to cc bcc message-id in-reply-to references)])\r\n")
    sockfd.sendStr("8 LOGOUT\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    sockfd.sendStr("9 UID FETCH 1 (BODY.PEEK[HEADER] BODY.PEEK[TEXT])\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    sockfd.sendStr("10 LOGOUT\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    time.sleep(10)
    #exit(0)
