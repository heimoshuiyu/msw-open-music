import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { CalcReadableFilesizeDetail } from "./Common";
import FfmpegConfig from "./FfmpegConfig";
import FileDialog from "./FileDialog";

function AudioPlayer(props) {
  // props.playingFile
  // props.setPlayingFile

  const [fileDialogShowStatus, setFileDialogShowStatus] = useState(false);
  const [loop, setLoop] = useState(true);
  const [raw, setRaw] = useState(false);
  const [prepare, setPrepare] = useState(false);
  const [selectedFfmpegConfig, setSelectedFfmpegConfig] = useState({});
  const [playingURL, setPlayingURL] = useState("");
  const [isPreparing, setIsPreparing] = useState(false);
  const [preparedFilesize, setPreparedFilesize] = useState(null);

  useEffect(() => {
    // no playing file
    if (props.playingFile.id === undefined) {
      setPlayingURL("");
      return;
    }
    if (raw) {
      console.log("Play raw file");
      setPlayingURL("/api/v1/get_file_direct?id=" + props.playingFile.id);
    } else {
      if (prepare) {
        // prepare file
        setIsPreparing(true);
        fetch("/api/v1/prepare_file_stream_direct", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            id: props.playingFile.id,
            config_name: selectedFfmpegConfig.name,
          }),
        })
          .then((response) => response.json())
          .then((data) => {
            setPreparedFilesize(data.filesize);
            setIsPreparing(false);
            setPlayingURL(
              `/api/v1/get_file_stream_direct?id=${props.playingFile.id}&config=${selectedFfmpegConfig.name}`
            );
          });
      } else {
        setPlayingURL(
          `/api/v1/get_file_stream?id=${props.playingFile.id}&config=${selectedFfmpegConfig.name}`
        );
      }
    }
  }, [props.playingFile.id, raw, prepare, selectedFfmpegConfig]);

  let navigate = useNavigate();

  return (
    <div>
      <h5>Player status</h5>
      {props.playingFile.id && (
        <span>
          <FileDialog
            showStatus={fileDialogShowStatus}
            setShowStatus={setFileDialogShowStatus}
            file={props.playingFile}
            setPlayingFile={() => {
              return;
            }}
          />

          <button
            onClick={() => {
              setFileDialogShowStatus(!fileDialogShowStatus);
            }}
          >
            {props.playingFile.filename}
          </button>

          <button
            onClick={() =>
              navigate(`/folders/${props.playingFile.folder_id}`)
            }
          >
            {props.playingFile.foldername}
          </button>

          <button disabled>
            {prepare
              ? CalcReadableFilesizeDetail(preparedFilesize)
              : CalcReadableFilesizeDetail(props.playingFile.filesize)}
          </button>

          {isPreparing && <button disabled>Preparing...</button>}

          {playingURL !== "" && (
            <button
              onClick={() => {
                props.setPlayingFile({});
              }}
            >
              Stop
            </button>
          )}
        </span>
      )}

      <br />

      <input
        checked={loop}
        onChange={(event) => setLoop(event.target.checked)}
        type="checkbox"
      />
      <label>Loop</label>

      <input
        checked={raw}
        onChange={(event) => setRaw(event.target.checked)}
        type="checkbox"
      />
      <label>Raw</label>

      {!raw && (
        <span>
          <input
            checked={prepare}
            onChange={(event) => setPrepare(event.target.checked)}
            type="checkbox"
          />
          <label>Prepare</label>
        </span>
      )}

      {playingURL !== "" && (
        <audio
          controls
          autoPlay
          loop={loop}
          className="audio-player"
          src={playingURL}
        ></audio>
      )}

      <FfmpegConfig
        selectedFfmpegConfig={selectedFfmpegConfig}
        setSelectedFfmpegConfig={setSelectedFfmpegConfig}
      />
    </div>
  );
}

export default AudioPlayer;
