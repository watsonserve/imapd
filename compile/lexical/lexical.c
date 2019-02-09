#include "lexical.h"

int endNumberConst(const char *str, int length)
{
    int i;

    if (!length)
    {
        length = strlen(str);
    }

    /* 10进制 */
    if ('0' != str[0])
    {
        for (i = 0; i < length && (('0' <= str[i] && str[i] <= '9') || '.' == str[i]); i++) {}
        return i;
    }

    /* hex */
    if ('x' == str[1] || 'X' == str[1])
    {
        for (i = 2; i < length && (('0' <= str[i] && str[i] <= '9') || ('A' <= str[i] && str[i] <= 'F')); i++) {}
        return i;
    }

    /* oct */
    for (i = 0; i < length && ('0' <= str[i] && str[i] < '8'); i++) {}
    return i;
}

int endStringConst(const char endCh, const char *str, int length)
{
    int i;

    if (!length)
    {
        length = strlen(str);
    }

    for (i = 1; i < length && endCh != str[i]; i++)
    {
        if('\\' == str[i] && (i + 1) < length)
        {
            i++;
        }
    }

    return i;
}

int endVariable(const char *str, int length)
{
    int i;

    if (!length)
    {
        length = strlen(str);
    }

    for (
        i = 0;
        (
            ('A' <= str[i] && str[i] <= 'Z') ||
            ('a' <= str[i] && str[i] <= 'z') ||
            ('0' <= str[i] && str[i] <= '9') ||
            '$' == str[i] ||
            '_' == str[i] ||
            0x80 <= str[i]
        );
        i++
    ) {}

    return i;
}

int isOperationSignal(const char ch)
{
    int i, length;

    const char lst[] = OPERATIONSIGNAL;
    length = sizeof lst;

    for (i = 0; i < length; i++)
    {
        if (ch == lst[i])
        {
            return 1;
        }
    }

    return 0;
}

const char * movPointer(const char *str, const int offset)
{
    return str + offset;
}
