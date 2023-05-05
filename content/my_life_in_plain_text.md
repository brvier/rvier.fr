Title: My life in plain text (🇬🇧)
Date: 2022-10-30
Modified: 2022-11-08
Category: Morg

For a long time i used Orgmode file to take notes, todos, and sometime agenda events. While i like a lot the Orgzly app on Android, i wasn't happy with the sync, and worse, without a specific tool the OrgMode format is quite unusable due to his complexity.

So i tryed many, many tools with Markdown support, but all seems missing somethi for my use:

- On Android :

	- Markor : No agenda
	- Joplin: No agenda and i don't like how the file are named (uuid.txt)
	- GitJournal : Only for journal

- On Linux :

	- There is todo.txt cli to manage todo in a plain text
	- vim to edit markdown

### Templates

My requirements are quite simple, so i opted for a simple files structure, and templates :

    todo.txt (todotxt format)
    done.txt (todotxt format)
    agenda.txt (agendatxt format)
    quicknote.txt (markdown)
    expenses.txt (expensetxt format)
    archives/
    attachments/
    journal/
      YYYY-mm-dd.txt
    notes/
      a_markdown_note.md
      an_other_markdown_note.md
      Projects/
        project_1.md
        project_2.md
    Posts/
      my_life_in_plain_text.md

#### Todotxt format

[Todo.txt: Future-proof task tracking in a file you control](http://todotxt.org/)


#### Agendatxt format

I ve created a agendatxt format similar to todotxt. Each line is an event, three format are accepted :

YYYY-mm-dd event for all day
YYYY-mm-dd HH:MM event at hour:minute
YYYY-mm-dd HH:MM HH:MM event with a start and a end

#### Expensetxt format

I ve created an expensetxt format similar to todotxt. Each line is an expense :

YYYY-mm-dd 000.00 short description

#### Markdown format

[Markdown — Wikipédia](https://fr.wikipedia.org/wiki/Markdown)

## On desktop

As a dev i use a lot vim, and terminals, i take notes with vim and use some aliases on zsh to quick add a todo or a journal entry.

## On mobile

This is more complex, as i need to visualize quickly events and todos, but also takes notes. So i write my own app. It s not finished yet, and look likes more a Proof of Concept than a finished app, many features are still missing, but i use it daily.

## Screenshot

![MOrg Screenshot](https://raw.githubusercontent.com/brvier/MOrg/master/screenshots/main_dark.jpg)

## Licence

And if you want to join me on this journey, you are more than welcomed, as it s distributed under MIT Licence.

[MOrg on GitHub](https://github.com/brvier/MOrg)
