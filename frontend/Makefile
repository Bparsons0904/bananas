.PHONY: install dev dev-react dev-solid dev-angular clean

install:
	cd react && npm install
	cd solid && npm install
	cd angular && npm install

dev:
	npx serve . -p 5172 --cors

dev-react:
	cd react && npm run dev

dev-solid:
	cd solid && npm run dev

dev-angular:
	cd angular && npm start

clean:
	rm -rf react/node_modules react/dist
	rm -rf solid/node_modules solid/dist
	rm -rf angular/node_modules angular/dist
