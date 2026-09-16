.PHONY: gpush

gpush:
	git add .
	git commit -m "${msg}"
	git push

