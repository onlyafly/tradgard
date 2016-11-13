- [x] pages served by name
- [x] translate page names back and forth between actual name and url
- [x] create page
  - [x] editing a non-existing page works
  - [x] viewing a non-existing page offers option to edit
- [x] link from page to page
- [x] link from root page to a page called "Home"
- [x] require password to log in
- [x] improve editing interface
  - [x] include markdown help
- [x] add simple styling
- [x] show nice error pages
- [x] provide some index of pages
- [ ] Add more robust error logging with levels
- [ ] https://github.com/gopherjs/gopherjs

# Viewing pages

- [x] Improve styling of blockquotes
- [ ] Link automatically back to all pages referencing this one

# Editing pages

- [x] Link to both basic markdown + blackfriday extensions + internal linking extension

## Markdown to HTML generation

- [ ] Fork and improve blackfriday
  - [ ] Generate headers with anchors
  - [ ] Generate table of contents automatically
  - [ ] Generate blockquotes with class="blockquote" to match bootstrap: http://v4-alpha.getbootstrap.com/content/typography/#blockquotes
  - [ ] External links open in new tab

## Markdown editor

- [ ] Unite the preview to match the exact output of the actual html generation
