'abcd' AND (SELECT 'a' FROM users WHERE username='administrator' AND LENGTH(password)>2)='a'