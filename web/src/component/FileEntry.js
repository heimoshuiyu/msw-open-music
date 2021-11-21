import { useState } from "react";
import { useNavigate } from "react-router";
import { CalcReadableFilesize } from "./Common";
import FileDialog from "./FileDialog";

function FileEntry(props) {
  const [showStatus, setShowStatus] = useState(false);
  let navigate = useNavigate();

  return (
    <tr>
      <td
        className="clickable"
        onClick={() => {
          // double click to play file and close dialog
          if (showStatus) {
            props.setPlayingFile(props.file);
            setShowStatus(false);
            return;
          }
          setShowStatus(true);
        }}
      >
        {props.file.filename}
      </td>
      <td
        className="clickable"
        onClick={() => navigate(`/search-folders/${props.file.folder_id}`)}
      >
        {props.file.foldername}
      </td>
      <td>
        {CalcReadableFilesize(props.file.filesize)}
        <FileDialog
          setPlayingFile={props.setPlayingFile}
          showStatus={showStatus}
          setShowStatus={setShowStatus}
          file={props.file}
        />
      </td>
    </tr>
  );
}

export default FileEntry;
