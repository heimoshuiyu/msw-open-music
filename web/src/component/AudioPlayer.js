import {useEffect, useState} from "react";
import {useNavigate} from "react-router";
import {CalcReadableFilesizeDetail} from "./Common";
import FfmpegConfig from "./FfmpegConfig";
import FileDialog from "./FileDialog";
import {Tr} from "../translate";

function AudioPlayer(props) {
  // props.playingFile
  // props.setPlayingFile

  const [fileDialogShowStatus, setFileDialogShowStatus] = useState(false);
  const [loop, setLoop] = useState(true);
  const [raw, setRaw] = useState(false);
  const [prepare, setPrepare] = useState(false);
  const [selectedFfmpegConfig, setSelectedFfmpegConfig] = useState({
    name: "",
    args: "",
  });
  const [playingURL, setPlayingURL] = useState("");
  const [isPreparing, setIsPreparing] = useState(false);
  const [timerCount, setTimerCount] = useState(0);
  const [timerID, setTimerID] = useState(null);
  const [beginPlayTime, setBeginPlayTime] = useState(null);
  const [lastID, setLastID] = useState(null);

  const recordPlaybackHistory = async (file_id, method, duration) => {
    await fetch('/api/v1/record_playback', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        playback: {
          file_id,
          method,
          duration,
        },
      })
    })
  }


  // init mediaSession API
  useEffect(() => {
    navigator.mediaSession.setActionHandler("stop", () => {
      props.setPlayingFile({});
    });
  }, []);

  useEffect(() => {
    // media session related staff
    navigator.mediaSession.metadata = new window.MediaMetadata({
      title: props.playingFile.filename,
      album: props.playingFile.foldername,
      artwork: [{src: "/favicon.png", type: "image/png"}],
    });
    // no playing file
    if (props.playingFile.id === undefined) {
      setPlayingURL("");
      return;
    }
    // crrently playing file, record interupt
    if (playingURL) {
      const endPlayTime = new Date()
      const duration = parseInt((endPlayTime - beginPlayTime) / 1000)
      recordPlaybackHistory(lastID, 2, duration)
    }
    setLastID(props.playingFile.id)
    // have playingFile, record begin time
    setBeginPlayTime(new Date())
    if (raw) {
      console.log("Play raw file");
      setPlayingURL("/api/v1/get_file_direct?id=" + props.playingFile.id);
    } else {
      if (prepare) {
        // prepare file
        setIsPreparing(true);
        fetch("/api/v1/prepare_file_stream_direct", {
          method: "POST",
          headers: {"Content-Type": "application/json"},
          body: JSON.stringify({
            id: props.playingFile.id,
            config_name: selectedFfmpegConfig.name,
          }),
        })
          .then((response) => response.json())
          .then((data) => {
            if (data.error) {
              alert(data.error);
              setIsPreparing(false);
              return;
            }
            props.setPlayingFile(data.file);
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
    <footer className="vertical">
      <h5>{Tr("Player status")}</h5>
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
            onClick={() => navigate(`/folders/${props.playingFile.folder_id}`)}
          >
            {props.playingFile.foldername}
          </button>

          <button disabled>
            {CalcReadableFilesizeDetail(props.playingFile.filesize)}
          </button>

          {isPreparing && <button disabled>{Tr("Preparing...")}</button>}

          {playingURL !== "" && (
            <button
              onClick={() => {
                props.setPlayingFile({});
                const endPlayTime = new Date()
                const duration = parseInt((endPlayTime - beginPlayTime) / 1000)
                setBeginPlayTime(endPlayTime)
                recordPlaybackHistory(props.playingFile.id, 3, duration)
              }}
            >
              {Tr("Stop")}
            </button>
          )}
        </span>
      )}

      <br />

      <span className="horizontal">
        <input
          className="number-input"
          disabled={timerID !== null}
          type="number"
          value={timerCount}
          onChange={(e) => {
            setTimerCount(e.target.value);
          }}
        />
        <button
          onClick={() => {
            if (timerID != null) {
              clearInterval(timerID);
              setTimerID(null);
              return;
            }
            setTimerID(
              setTimeout(() => {
                props.setPlayingFile({});
                setTimerID(null);
              }, timerCount * 1000 * 60)
            );
          }}
        >
          {Tr("Stop Timer")}
        </button>
      </span>

      <span>
        <span>
          <input
            checked={loop}
            onChange={(event) => setLoop(event.target.checked)}
            type="checkbox"
          />
          <label>{Tr("Loop")}</label>
        </span>

        <span>
          <input
            checked={raw}
            onChange={(event) => setRaw(event.target.checked)}
            type="checkbox"
          />
          <label>{Tr("Raw")}</label>
        </span>

        {!raw && (
          <span>
            <input
              checked={prepare}
              onChange={(event) => setPrepare(event.target.checked)}
              type="checkbox"
            />
            <label>{Tr("Prepare")}</label>
          </span>
        )}
      </span>

      <audio
        id="dom-player"
        controls
        autoPlay
        className="audio-player"
        src={playingURL}
        onEnded={async () => {
          const player = document.getElementById('dom-player')
          if (loop) {
            player.play()
          }
          const endPlayTime = new Date()
          const duration = parseInt((endPlayTime - beginPlayTime) / 1000)
          setBeginPlayTime(endPlayTime)
          recordPlaybackHistory(props.playingFile.id, 1, duration)
        }}
      ></audio>

      <FfmpegConfig
        selectedFfmpegConfig={selectedFfmpegConfig}
        setSelectedFfmpegConfig={setSelectedFfmpegConfig}
      />
    </footer >
  );
}

export default AudioPlayer;
