import os

if __name__ == '__main__':
    html = open('index.html').read()
    html = html.replace('MAPBOX_TOKEN', os.environ['MAPBOX_TOKEN'])
    open('index.html', 'w+').write(html)
