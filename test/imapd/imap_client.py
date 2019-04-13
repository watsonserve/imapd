# coding=utf-8

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
    if '2 NO LOGIN failed.' == buf[0:18]:
        exit(0)

    sendStr("3 CAPABILITY\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    sendStr("4 SELECT Inbox\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    sendStr("5 LIST \"\" \"Inbox/*\"\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    # 检索有多少邮件，并查看当前邮箱中所有邮件的UID
    # UID SEARCH 中 数字表示当前邮箱中的第offset个消息：第offset个消息，例如2:4为第2条消息到第4条消息共3个
    sendStr("6 UID SEARCH 2:4 NOT DELETED\r\n")
    # sendStr("6 UID SEARCH 1:* NOT DELETED\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    # sendStr("7 UID FETCH 12365:12372 (INTERNALDATE UID RFC822.SIZE FLAGS BODY.PEEK[HEADER.FIELDS (date subject from content-type to cc bcc message-id in-reply-to references)])\r\n")
    # buf = sockfd.recv(BUFSIZ)
    # printB(buf)

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
