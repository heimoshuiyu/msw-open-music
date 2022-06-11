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
      <p
        style={{
          cursor: "pointer",
        }}
        onClick={() => {
          props.setPlayingFile(props.file);
          props.setShowStatus(false);
        }}
      >
        {props.file.filename}
      </p>
      <p>
        Play: play using browser player.
        <br />
        Info for more actions.
      </p>
      <button
        onClick={() => {
          navigate(`/files/${props.file.id}`);
          props.setShowStatus(false);
        }}
      >
        Info
      </button>
      <button
        onClick={() => {
          props.setPlayingFile(props.file);
          props.setShowStatus(false);
        }}
      >
        Play
      </button>
      <button onClick={() => props.setShowStatus(false)}>Close</button>
    </dialog>
  );
}

export default FileDialog;
