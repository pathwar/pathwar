const fs = require('fs'),
      glob = require('glob'),
      path = require('path'),
      mime = require('mime-types');

const rootDir = path.join(__dirname, '../../');
const contentDir = path.join(__dirname, '../content');
const assetsDir = path.join(__dirname, '../assets/files');

[contentDir, assetsDir].forEach(dir => {
  if (!fs.existsSync(dir)){
    fs.mkdirSync(dir, { recursive: true });
  }

  const files = fs.readdirSync(dir);
  for (const file of files) {
    fs.unlink(path.join(dir, file), err => {
      if (err) throw err;
    });
  }
})

const options = {
  cwd: rootDir,
  ignore: []
};

const folders = glob.sync('**/', options);

folders.forEach(folder => {

  folder = folder
    .replace(/.$/, '');

  const files = fs
    .readdirSync(path.join(rootDir, folder), { withFileTypes: true })
    .filter(item => !item.isDirectory())
    .map(file => file.name)

  const mainFile = files
    .find(file => file.startsWith('1_'));

  const folderId = folder
    .replace(/\/+/g, '-');

  const category = folder
    .split('/')[0];

  const filesToShowinList = mainFile
    ? [mainFile]
    : files;

  let parentFolder = folder.split('/');
      parentFolder = parentFolder[parentFolder.length - 1];

  const parentFolderTitle = parentFolder
    .replace(/\_+/g, ' ');

  files.forEach(file => {
    const title = path.basename(file);
    const filename = folder
      .replace(`${category}/`, '')
      .replace(/\/+/g, '-')
      .concat(`--${file}`);

    let mediaType = mime.lookup(file);

    // Handle unrecognized extensions, else stop
    if (!mediaType) {
      if (file.endsWith('.sketch')) {
        mediaType = 'application/sketch';
      } else {
        return;
      }
    }

    if (file.includes('zip')) {
      console.log(file)
    }

    const mediaTypeSplit = mediaType.split('/'),
          mediaMainType = mediaTypeSplit[0],
          mediaSubType = mediaTypeSplit[1];

    let fm = "---\n";
    fm += `title: ${title}\n`;
    fm += `show_in_list: ${filesToShowinList.includes(file)}\n`;
    fm += `folder_id: ${folderId}\n`;
    fm += `categories: ["${category}"]\n`;
    fm += `media_type: ${mediaType}\n`;
    fm += `media_type_main: ${mediaMainType}\n`;
    fm += `media_type_sub: ${mediaSubType}\n`;
    fm += `parent_folder: ${parentFolder}\n`;
    fm += `parent_folder_title: ${parentFolderTitle}\n`;
    fm += `file_path: /files/${filename}\n`;
    fm += `file_name: ${file}\n`;
    fm += '---';

    fs.writeFileSync(path.join(contentDir, `${filename}.md`), fm);
    fs.copyFileSync(path.join(rootDir, folder, file), path.join(assetsDir, filename));
  })
});
