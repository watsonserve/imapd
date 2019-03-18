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

def tcp(server_name, port, use_ssl):
    sockfd = socket.socket(socket.AF_INET)
    if use_ssl:
        # context = ssl._create_unverified_context()
        context = ssl.create_default_context()
        sockfd = context.wrap_socket(sockfd, server_hostname=server_name)
    sockfd.connect((server_name, port))
    return sockfd

if "__main__" == __name__:
    BUFSIZ = 4096
    # sockfd = tcp("imap-mail.outlook.com", 993, True)
    sockfd = tcp("imap.watsonserve.com", 993, True)
    sendStr = foo(sockfd)
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    print("connected!")
    sendStr("1 CAPABILITY\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    # sendStr("2 AUTHENTICATE PLAIN\r\n")
    sendStr("2 LOGIN james@watsonserve.com 123456\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    # sendStr("*\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)

    sendStr("3 CAPABILITY\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    sendStr("4 ID (\"name\" \"py mail\" \"version\" \"0.0.1\")\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    # sendStr("5 SELECT Notes\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)

    # sendStr("6 LIST \"\" \"\"\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)

    sendStr("6 LIST \"INBOX\" \"*\"\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    # sendStr("7 LIST \"*\" \"\"\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)

    # sendStr("4 SELECT DELETED\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)

    # sendStr("5 UID SEARCH 143:* NOT DELETED\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)
    #sendStr("6 UID FETCH 126:127 (INTERNALDATE UID RFC822.SIZE FLAGS BODY.PEEK[HEADER.FIELDS (date subject from content-type to cc bcc message-id in-reply-to references)])\r\n")
    #buf = sockfd.recv(BUFSIZ)
    #printB(buf)
    # sendStr("6 NOOP\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)
    # sendStr("7 UID FETCH 126 (BODYSTRUCTURE BODY.PEEK[HEADER])\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)
    # sendStr("8 UID FETCH 1:3 (INTERNALDATE UID RFC822.SIZE FLAGS BODY.PEEK[HEADER.FIELDS (date subject from content-type to cc bcc message-id in-reply-to references)])\r\n")
    # sendStr("8 LOGOUT\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)
    # sendStr("9 UID FETCH 1 (BODY.PEEK[HEADER] BODY.PEEK[TEXT])\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)
    sendStr("10 LOGOUT\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    time.sleep(1)
    #exit(0)
