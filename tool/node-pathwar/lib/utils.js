module.exports.getVersion = function(module) {
  return require(module + '/package.json').version;
};
