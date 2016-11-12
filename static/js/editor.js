// Adapted from https://github.com/NextStepWebs/simplemde-markdown-editor/blob/master/src/js/simplemde.js, line 1069
var fullToolbar = [
	{
		name: "bold",
		action: SimpleMDE.toggleBold,
		className: "fa fa-bold",
		title: "Bold"
	},
	{
		name: "italic",
		action: SimpleMDE.toggleItalic,
		className: "fa fa-italic",
		title: "Italic"
	},
	{
		name: "strikethrough",
		action: SimpleMDE.toggleStrikethrough,
		className: "fa fa-strikethrough",
		title: "Strikethrough"
	},
	{
		name: "heading",
		action: SimpleMDE.toggleHeadingSmaller,
		className: "fa fa-header",
		title: "Heading"
	},
	{
		name: "heading-smaller",
		action: SimpleMDE.toggleHeadingSmaller,
		className: "fa fa-header fa-header-x fa-header-smaller",
		title: "Smaller Heading"
	},
	{
		name: "heading-bigger",
		action: SimpleMDE.toggleHeadingBigger,
		className: "fa fa-header fa-header-x fa-header-bigger",
		title: "Bigger Heading"
	},
	{
		name: "heading-1",
		action: SimpleMDE.toggleHeading1,
		className: "fa fa-header fa-header-x fa-header-1",
		title: "Big Heading"
	},
	{
		name: "heading-2",
		action: SimpleMDE.toggleHeading2,
		className: "fa fa-header fa-header-x fa-header-2",
		title: "Medium Heading"
	},
	{
		name: "heading-3",
		action: SimpleMDE.toggleHeading3,
		className: "fa fa-header fa-header-x fa-header-3",
		title: "Small Heading"
	},
	"|",
	{
		name: "code",
		action: SimpleMDE.toggleCodeBlock,
		className: "fa fa-code",
		title: "Code"
	},
	{
		name: "quote",
		action: SimpleMDE.toggleBlockquote,
		className: "fa fa-quote-left",
		title: "Quote"
	},
	{
		name: "unordered-list",
		action: SimpleMDE.toggleUnorderedList,
		className: "fa fa-list-ul",
		title: "Generic List"
	},
	{
		name: "ordered-list",
		action: SimpleMDE.toggleOrderedList,
		className: "fa fa-list-ol",
		title: "Numbered List"
	},
	{
		name: "clean-block",
		action: SimpleMDE.cleanBlock,
		className: "fa fa-eraser fa-clean-block",
		title: "Clean block"
	},
	"|",
	{
		name: "link",
		action: SimpleMDE.drawLink,
		className: "fa fa-link",
		title: "Create Link"
	},
	{
		name: "image",
		action: SimpleMDE.drawImage,
		className: "fa fa-picture-o",
		title: "Insert Image"
	},
	{
		name: "table",
		action: SimpleMDE.drawTable,
		className: "fa fa-table",
		title: "Insert Table"
	},
	{
		name: "horizontal-rule",
		action: SimpleMDE.drawHorizontalRule,
		className: "fa fa-minus",
		title: "Insert Horizontal Line"
	},
	"|",
	{
		name: "preview",
		action: SimpleMDE.togglePreview,
		className: "fa fa-eye no-disable",
		title: "Toggle Preview"
	},
	{
		name: "side-by-side",
		action: SimpleMDE.toggleSideBySide,
		className: "fa fa-columns no-disable no-mobile",
		title: "Toggle Side by Side"
	},
	{
		name: "fullscreen",
		action: SimpleMDE.toggleFullScreen,
		className: "fa fa-arrows-alt no-disable no-mobile",
		title: "Toggle Fullscreen"
	},
	"|",
	{
		name: "guide",
		action: "https://simplemde.com/markdown-guide",
		className: "fa fa-question-circle",
		title: "Markdown Guide"
	},
	"|",
	{
		name: "undo",
		action: SimpleMDE.undo,
		className: "fa fa-undo no-disable",
		title: "Undo"
	},
	{
		name: "redo",
		action: SimpleMDE.redo,
		className: "fa fa-repeat no-disable",
		title: "Redo"
	}
];

function codeMirrorReplaceSelection(cm, startEnd) {
	if(/editor-preview-active/.test(cm.getWrapperElement().lastChild.className))
		return;

	var text;
	var start = startEnd[0];
	var end = startEnd[1];
	var startPoint = cm.getCursor("start");
	var endPoint = cm.getCursor("end");

	text = cm.getSelection();
	cm.replaceSelection(start + text + end);

	startPoint.ch += start.length;
	if(startPoint !== endPoint) {
		endPoint.ch += start.length;
	}

	cm.setSelection(startPoint, endPoint);
	cm.focus();
}

/**
 * Action for drawing a page link.
 */
function drawPageLink(editor) {
	codeMirrorReplaceSelection(editor.codemirror, editor.options.insertTexts.pageLink);
}

var customToolbar = [
	{
		name: "bold",
		action: SimpleMDE.toggleBold,
		className: "fa fa-bold",
		title: "Bold"
	},
	{
		name: "italic",
		action: SimpleMDE.toggleItalic,
		className: "fa fa-italic",
		title: "Italic"
	},
	{
		name: "strikethrough",
		action: SimpleMDE.toggleStrikethrough,
		className: "fa fa-strikethrough",
		title: "Strikethrough"
	},
  "|",
	{
		name: "heading-1",
		action: SimpleMDE.toggleHeading1,
		className: "fa fa-header fa-header-x fa-header-1",
		title: "Big Heading"
	},
	{
		name: "heading-2",
		action: SimpleMDE.toggleHeading2,
		className: "fa fa-header fa-header-x fa-header-2",
		title: "Medium Heading"
	},
	{
		name: "heading-3",
		action: SimpleMDE.toggleHeading3,
		className: "fa fa-header fa-header-x fa-header-3",
		title: "Small Heading"
	},
	"|",
  {
    name: "quote",
    action: SimpleMDE.toggleBlockquote,
    className: "fa fa-quote-left",
    title: "Quote"
  },
  {
    name: "unordered-list",
    action: SimpleMDE.toggleUnorderedList,
    className: "fa fa-list-ul",
    title: "Generic List"
  },
  {
    name: "ordered-list",
    action: SimpleMDE.toggleOrderedList,
    className: "fa fa-list-ol",
    title: "Numbered List"
  },
	"|",
	{
		name: "link",
		action: SimpleMDE.drawLink,
		className: "fa fa-external-link",
		title: "External Link"
	},
  {
    name: "page",
    action: drawPageLink,
    className: "fa fa-link",
    title: "Link to Page",
  },
	{
		name: "image",
		action: SimpleMDE.drawImage,
		className: "fa fa-picture-o",
		title: "Insert Image"
	},
	{
		name: "table",
		action: SimpleMDE.drawTable,
		className: "fa fa-table",
		title: "Insert Table"
	},
	"|",
	{
		name: "preview",
		action: SimpleMDE.togglePreview,
		className: "fa fa-eye no-disable",
		title: "Toggle Preview"
	},
	{
		name: "guide",
		action: "https://simplemde.com/markdown-guide",
		className: "fa fa-question-circle",
		title: "Markdown Guide"
	}
];

// https://github.com/NextStepWebs/simplemde-markdown-editor
var simplemde = new SimpleMDE({
  element: document.getElementById("pageContentEditor"),
  toolbar: customToolbar,
  forceSync: true,
  insertTexts: {
    horizontalRule: ["", "\n\n-----\n\n"],
    image: ["![](http://", ")"],
    link: ["[", "](http://)"],
    pageLink: ["{", "}"],
    table: ["", "\n\n| Column 1 | Column 2 | Column 3 |\n| -------- | -------- | -------- |\n| Text     | Text      | Text     |\n\n"],
  }
});
