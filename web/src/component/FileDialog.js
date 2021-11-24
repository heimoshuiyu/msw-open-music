import { useNavigate } from "react-router";

function FileDialog(props) {
  // props.showStatus
  // props.setShowStatus
  // props.playingFile
  // props.setPlayingFile
  // props.file

  let navigate = useNavigate();

  return (
    <dialog open={props.showStatus}>
      <p>{props.file.filename}</p>
      <p>
        Download 使用 Axios 异步下载
        <br />
        Play 调用网页播放器播放
        <br />
      </p>
      <button>Download</button>
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
          navigate(`/share/${props.file.id}`);
          props.setShowStatus(false);
        }}
      >
        Share
      </button>
      <button onClick={() => props.setShowStatus(false)}>Close</button>
    </dialog>
  );
}

export default FileDialog;
