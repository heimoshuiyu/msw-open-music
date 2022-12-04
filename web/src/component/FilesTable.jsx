import * as React from 'react';
import FileEntry from "./FileEntry";
import { Tr } from "../translate";

function FilesTable(props) {
  if (props.files.length === 0) {
    return null;
  }
  return (
    <table>
      <thead>
        <tr>
          <th>{Tr("Filename")}</th>
          <th>{Tr("Folder Name")}</th>
          <th>{Tr("Size")}</th>
        </tr>
      </thead>
      <tbody>
        {props.files.map((file) => (
          <FileEntry
            setPlayingFile={props.setPlayingFile}
            key={file.id}
            file={file}
          />
        ))}
      </tbody>
    </table>
  );
}

export default FilesTable;
