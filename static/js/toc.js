var headerToIdent = function(input) {
  return input.replace(/[^a-zA-Z]/g, "_").toLowerCase();
};

var nav =
  "<nav role='navigation' class='table-of-contents'>" +
  "<h2>Contents</h2>" +
  "<ul>";

var processHeader = function() {
  var el = $(this);
  var title = el.text();
  var ident = headerToIdent(title);
  el.attr("id", ident);
  var link = "#" + ident;

  var newItem =
    "<li>" +
      "<a href='" + link + "'>" +
        title +
      "</a>" +
    "</li>";

  nav += newItem;
};

$("article h1").each(processHeader);
$("article h2").each(processHeader);
$("article h3").each(processHeader);
$("article h4").each(processHeader);
$("article h5").each(processHeader);

nav +=
  "</ul>" +
  "</nav>";

$("article").prepend(nav);
