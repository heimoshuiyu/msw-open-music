import { useNavigate } from "react-router";

function FileDialog(props) {
  // props.showStatus
  // props.setShowStatus
  // props.playingFile
  // props.setPlayingFile
  // props.file

  let navigate = useNavigate();

  const downloadURL = "/api/v1/get_file_direct?id=" + props.file.id;

  return (
    <dialog open={props.showStatus}>
      <p>{props.file.filename}</p>
      <p>
        Download 使用浏览器下载原文件
        <br />
        Play 调用网页播放器播放
        <br />
      </p>
      <a href={downloadURL} download>
        <button>Download</button>
      </a>
      <button
        onClick={() => {
          props.setPlayingFile(props.file);
          props.setShowStatus(false);
        }}
      >
        Play
      </button>
      <button
        onClick={() => {
          navigate(`/files/${props.file.id}`);
          props.setShowStatus(false);
        }}
      >
        Info
      </button>
      <button onClick={() => props.setShowStatus(false)}>Close</button>
    </dialog>
  );
}

export default FileDialog;
