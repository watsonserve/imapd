#ifndef _LEXICAL_H_
#define _LEXICAL_H_

#include <string.h>

#define OPERATIONSIGNAL ",.(){}[]+-*/%=&|~^!><;?:@#`\\\n"

const char * movPointer(const char *, const int);

int endNumberConst(const char *, int);

int endStringConst(const char, const char *, int);

int endVariable(const char *, int);

int isOperationSignal(char);

#endif
