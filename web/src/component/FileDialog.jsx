import * as React from 'react';
import { useNavigate } from "react-router";
import { Tr } from "../translate";

function FileDialog(props) {
  // props.showStatus
  // props.setShowStatus
  // props.playingFile
  // props.setPlayingFile
  // props.file

  let navigate = useNavigate();

  return (
    <dialog
      open={props.showStatus}
      style={{
        zIndex: 1,
      }}
    >
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
        {Tr("Play: play using browser player.")}
        <br />
        {Tr("Info for more actions.")}
      </p>
      <button
        onClick={() => {
          navigate(`/files/${props.file.id}`);
          props.setShowStatus(false);
        }}
      >
        {Tr("Info")}
      </button>
      <button
        onClick={() => {
          props.setPlayingFile(props.file);
          props.setShowStatus(false);
        }}
      >
        {Tr("Play")}
      </button>
      <button onClick={() => props.setShowStatus(false)}>{Tr("Close")}</button>
    </dialog>
  );
}

export default FileDialog;
