var headerToIdent = function(input) {
  return input.replace(/[^a-zA-Z]/g, "_").toLowerCase();
};

var output = "";

var processHeader = function() {
  var el = $(this);
  var tag = el.prop("tagName").toLowerCase();
  var title = el.text();
  var ident = headerToIdent(title);
  el.attr("id", ident);
  var link = "#" + ident;

  var newItem =
    "<li class='toc-li-" + tag + "'>" +
      "<a href='" + link + "'>" +
        title +
      "</a>" +
    "</li>";

  output += newItem;
};

$("article h1, article h2, article h3, article h4, article h5, article h6").each(processHeader);

if (output.length > 0) {
  var nav =
    "<nav role='navigation' class='table-of-contents'>" +
    "<h2>Contents</h2>" +
    "<ul>" +
    output +
    "</ul>" +
    "</nav>";

  $("article").prepend(nav);
}
