import json
import requests
from bs4 import BeautifulSoup

symbols = {}

for ch in "abcdefghijklmnopqrstuvwxyz":
    char_symbols = []

    url = f"https://www.advfn.com/nyse/newyorkstockexchange.asp?companies={ch.upper()}"
    print("requesting", url, "...")

    resp = requests.get(url)
    soup = BeautifulSoup(resp.content, features='lxml')

    market_table = soup.find('table', class_='market')
    for row in market_table.find_all('tr'):
        td = row.find_all('td')
        if len(td) != 3:
            continue
        
        # read data
        name = td[0].string.strip()
        symbol = td[1].string.strip()

        if symbol in symbols:
            print("")
            print("WARN ::", symbol, "already in symbols!")
            continue

        symbols[symbol] = name

    print("loaded:", len(symbols), "symbols")

print(json.dumps(symbols))
