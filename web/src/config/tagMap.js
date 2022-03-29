const tags = [
  {
    name: "web",
    color: "#a2f6a9",
  },
  {
    name: "intro",
    color: "#92c4e9",
  },
  {
    name: "tutorial",
    color: "#b6e9d8",
  },
  {
    name: "sql",
    color: "#ead597",
  },
  {
    name: "injection",
    color: "#a6abff",
  },
  {
    name: "bruteforce",
    color: "#e79cfd",
  },
  {
    name: "exec",
    color: "#a59ecb",
  },
  {
    name: "http",
    color: "#c8fef6",
  },
  {
    name: "data-leak",
    color: "#8395bc",
  },
];

export const getTagColor = tag => {
  return tags.find(tagObj => tagObj.name === tag).color;
};

export default tags;
