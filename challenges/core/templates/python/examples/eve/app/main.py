# -*- coding: utf-8 -*-

from eve import Eve

if __name__ == '__main__':
    app = Eve()
    app.run(host='0.0.0.0', port=8080)
