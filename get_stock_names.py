import requests
from bs4 import BeautifulSoup

symbols = []
for ch in "abcdefghijklmnopqrstuvwxyz":
    url = f"https://www.advfn.com/nasdaq/nasdaq.asp?companies={ch.upper()}"
    print("requesting", url, "...")

    resp = requests.get(url)
    soup = BeautifulSoup(resp.content, features='lxml')

    market_table = soup.find('table', class_='market')
    char_symbols = []
    for row in market_table.find_all('tr'):
        td = row.find_all('td')
        if len(td) != 3:
            continue
        symbol = td[1].string
        symbols.append(symbol)
        char_symbols.append(symbol)
    print()
    print("+ loaded", len(char_symbols), "for char", ch.upper())
    print("- that's a total of", len(symbols), "symbols!")

print("symbols: (", len(symbols), "):")
print(symbols)
