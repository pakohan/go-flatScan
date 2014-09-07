(function() {
  var template = Handlebars.template, templates = Handlebars.templates = Handlebars.templates || {};
templates['table.template'] = template({"1":function(depth0,helpers,partials,data) {
  var stack1, helper, helperMissing=helpers.helperMissing, escapeExpression=this.escapeExpression, functionType="function", buffer = "\n    <tr>\n      <td>"
    + escapeExpression((helper = helpers.unixtime || (depth0 && depth0.unixtime) || helperMissing,helper.call(depth0, (depth0 && depth0.TimeUpdated), {"name":"unixtime","hash":{},"data":data})))
    + "</td>\n      <td>"
    + escapeExpression(((helper = helpers.RentN || (depth0 && depth0.RentN)),(typeof helper === functionType ? helper.call(depth0, {"name":"RentN","hash":{},"data":data}) : helper)))
    + " &euro;</td>\n      <td>\n        <button class=\"btn btn-primary btn-sm\" data-toggle=\"modal\" data-target=\"#myModal-"
    + escapeExpression(((stack1 = (data == null || data === false ? data : data.index)),typeof stack1 === functionType ? stack1.apply(depth0) : stack1))
    + "\">"
    + escapeExpression(((helper = helpers.Title || (depth0 && depth0.Title)),(typeof helper === functionType ? helper.call(depth0, {"name":"Title","hash":{},"data":data}) : helper)))
    + "</button>\n      </td>\n      <td>"
    + escapeExpression(((helper = helpers.Street || (depth0 && depth0.Street)),(typeof helper === functionType ? helper.call(depth0, {"name":"Street","hash":{},"data":data}) : helper)))
    + " "
    + escapeExpression(((helper = helpers.Zip || (depth0 && depth0.Zip)),(typeof helper === functionType ? helper.call(depth0, {"name":"Zip","hash":{},"data":data}) : helper)))
    + " "
    + escapeExpression(((helper = helpers.District || (depth0 && depth0.District)),(typeof helper === functionType ? helper.call(depth0, {"name":"District","hash":{},"data":data}) : helper)))
    + "</td>\n      <td>"
    + escapeExpression(((helper = helpers.Rooms || (depth0 && depth0.Rooms)),(typeof helper === functionType ? helper.call(depth0, {"name":"Rooms","hash":{},"data":data}) : helper)))
    + "</td>\n      <td>"
    + escapeExpression(((helper = helpers.Size || (depth0 && depth0.Size)),(typeof helper === functionType ? helper.call(depth0, {"name":"Size","hash":{},"data":data}) : helper)))
    + "</td>\n      <td><a href=\"http://kleinanzeigen.ebay.de"
    + escapeExpression(((helper = helpers.Url || (depth0 && depth0.Url)),(typeof helper === functionType ? helper.call(depth0, {"name":"Url","hash":{},"data":data}) : helper)))
    + "\">Angebot</a></td>\n    </tr>\n    <div class=\"modal fade\" id=\"myModal-"
    + escapeExpression(((stack1 = (data == null || data === false ? data : data.index)),typeof stack1 === functionType ? stack1.apply(depth0) : stack1))
    + "\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"modal-label\" aria-hidden=\"true\">\n      <div class=\"modal-dialog\">\n        <div class=\"modal-content\">\n          <div class=\"modal-header\">\n            <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-hidden=\"true\">&times;</button>\n            <h4 class=\"modal-title\" id=\"modal-label\">"
    + escapeExpression(((helper = helpers.Title || (depth0 && depth0.Title)),(typeof helper === functionType ? helper.call(depth0, {"name":"Title","hash":{},"data":data}) : helper)))
    + "</h4>\n          </div>\n          <div class=\"modal-body\">"
    + escapeExpression(((helper = helpers.Street || (depth0 && depth0.Street)),(typeof helper === functionType ? helper.call(depth0, {"name":"Street","hash":{},"data":data}) : helper)))
    + " "
    + escapeExpression(((helper = helpers.Zip || (depth0 && depth0.Zip)),(typeof helper === functionType ? helper.call(depth0, {"name":"Zip","hash":{},"data":data}) : helper)))
    + " "
    + escapeExpression(((helper = helpers.District || (depth0 && depth0.District)),(typeof helper === functionType ? helper.call(depth0, {"name":"District","hash":{},"data":data}) : helper)))
    + "<br>\n            <img src=\"//maps.googleapis.com/maps/api/staticmap?center="
    + escapeExpression(((helper = helpers.Street || (depth0 && depth0.Street)),(typeof helper === functionType ? helper.call(depth0, {"name":"Street","hash":{},"data":data}) : helper)))
    + ","
    + escapeExpression(((helper = helpers.Zip || (depth0 && depth0.Zip)),(typeof helper === functionType ? helper.call(depth0, {"name":"Zip","hash":{},"data":data}) : helper)))
    + ","
    + escapeExpression(((helper = helpers.District || (depth0 && depth0.District)),(typeof helper === functionType ? helper.call(depth0, {"name":"District","hash":{},"data":data}) : helper)))
    + "&zoom=13&size=560x560&sensor=false&key=AIzaSyBqfaRyN-jIdyEVU7RtUUF2Q1sCVMk-SOE\"><br>\n            ";
  stack1 = (helper = helpers.replace || (depth0 && depth0.replace) || helperMissing,helper.call(depth0, (depth0 && depth0.Description), {"name":"replace","hash":{},"data":data}));
  if(stack1 || stack1 === 0) { buffer += stack1; }
  return buffer + "\n          </div>\n          <div class=\"modal-footer\">\n            <button type=\"button\" class=\"btn btn-default\" data-dismiss=\"modal\">Close</button>\n          </div>\n        </div>\n      </div>\n    </div>\n    ";
},"compiler":[5,">= 2.0.0"],"main":function(depth0,helpers,partials,data) {
  var stack1, helper, options, functionType="function", blockHelperMissing=helpers.blockHelperMissing, buffer = "<table cellpadding=\"0\" cellspacing=\"0\" border=\"0\" class=\"datatable table table-striped\">\n  <thead>\n    <th>Hinzugef&uuml;gt</th>\n    <th>Miete</th>\n    <th>Überschrift</th>\n    <th>Adresse</th>\n    <th>R&auml;ume</th>\n    <th>Fl&auml;che</th>\n    <th>Link</th>\n  </thead>\n  <tbody>\n    ";
  stack1 = ((helper = helpers.offers || (depth0 && depth0.offers)),(options={"name":"offers","hash":{},"fn":this.program(1, data),"inverse":this.noop,"data":data}),(typeof helper === functionType ? helper.call(depth0, options) : helper));
  if (!helpers.offers) { stack1 = blockHelperMissing.call(depth0, stack1, options); }
  if(stack1 || stack1 === 0) { buffer += stack1; }
  return buffer + "\n  </tbody>\n</table>\n";
},"useData":true});
})();