Title: Publish your Android Kivy app to F-droid
Date: 2021-11-11
Status: draft
Lang: EN

It s quite easy to publish your Kivy app to F-Droid.

- Fork on gitlab `https://gitlab.com/fdroid/fdroiddata/`
- create a new branch `td.domain.yourapp`
- install f-droid tools `pip install git+https://gitlab.com/fdroid/fdroidserver.git`
- add a `td.domain.yourappname.yml` in metadata with the following template :
```
	Categories:
	  - Writing
	License: MIT
	WebSite: https://rvier.fr/pages/morg.html
	SourceCode: https://github.com/brvier/morg
	IssueTracker: https://github.com/brvier/morg/issues
	Changelog: https://raw.githubusercontent.com/brvier/MOrg/HEAD/CHANGELOG.md

	Name: MOrg
	Description: |-
		A future proof opinionated software to manage your life in plaintext : todo, agenda, journal and notes.
		Work In Progress : many features still missing

	RepoType: git
	Repo: https://github.com/brvier/morg

	Builds:
	  - versionName: 0.2.1
		versionCode: 8193
		commit: f990b1003a1e2a92040e15a969effc3ea83e169f
		timeout: 10800
		sudo:
		  - apt-get update
		  - sudo apt install -y git zip unzip openjdk-17-jdk python3-pip autoconf libtool
			pkg-config zlib1g-dev libncurses5-dev libncursesw5-dev libtinfo5 cmake libffi-dev
			libssl-dev python3 make build-essential libltdl-dev autoconf automake ccache
			gettext patch python3-pip python3-setuptools sudo unzip zip
		  - pip3 install --upgrade Cython==0.29.19 virtualenv
		  - pip3 install --upgrade kivy
		  - pip3 install --upgrade buildozer
		output: bin/morg-$$VERSION$$-arm64-v8a_armeabi-v7a-release.aab
		build: buildozer android release
		ndk: r23b

	AutoUpdateMode: None
	UpdateCheckMode: None
	CurrentVersion: 0.2.1
	CurrentVersionCode: 8193
```
- test your yaml metadata : `fdroid readmeta && fdroid lint td.domain.yourappname`
- git commit and merge
- see the CI/CD build
- if succesfull create a merge request and follow the template

and now wait until merged :)
