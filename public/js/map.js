window.ZipAreasMap || (window.ZipAreasMap = {});
ZipAreasMap.zipCodes = zips

ZipAreasMap.initialize = function() {
  var canvas, map, mapOptions;
  canvas = $('#zip-area-map-canvas').get(0);
  if (canvas) {
    mapOptions = {
      center: new google.maps.LatLng(52.5214359,13.4057142),
      zoom: 12
    };
    map = new google.maps.Map(canvas, mapOptions);
    ZipAreasMap.polygons = [];
    ZipAreasMap.zipCodes.forEach(function(zipArea) {
      return ZipAreasMap.polygons.push(zipArea.polygon());
    });
    ZipAreasMap.setActiveZipCodes = function() {
      var activePolygons;
      activePolygons = ZipAreasMap.polygons.filter(function(polygon) {
        return polygon.active === true;
      });
      ZipAreasMap.selectedZipCodes = activePolygons.map(function(polygon) {
        return polygon.zipCode;
      });
      return $('#selected-zip-codes').val(ZipAreasMap.selectedZipCodes.join(';'));
    };
    return ZipAreasMap.polygons.forEach(function(polygon) {
      polygon.setMap(map);
      return google.maps.event.addListener(polygon, "click", function() {
        var fillOpacity;
        fillOpacity = (this.fillOpacity === 0 ? 0.5 : 0);
        this.setOptions({
          fillOpacity: fillOpacity
        });
        this.active = !this.active;
        return ZipAreasMap.setActiveZipCodes();
      });
    });
  }
};

$(function() {
  ZipAreasMap.selectedZipCodes = $('#selected-zip-codes').val().split(';');
  ZipAreasMap.initialize();
  $('#selected-zip-codes').text(ZipAreasMap.selectedZipCodes);
  return $('#change-color').click(function() {
    var color;
    color = RandomColor.generate();
    return ZipAreasMap.polygons.forEach(function(polygon) {
      return polygon.setOptions({
        fillColor: color,
        strokeColor: color
      });
    });
  });
});