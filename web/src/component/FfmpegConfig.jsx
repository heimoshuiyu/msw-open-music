import * as React from 'react';
import { useEffect, useState } from "react";

function FfmpegConfig(props) {
  // props.setSelectedFfmpegConfig
  // props.selectedFfmpegConfig

  const [ffmpegConfigList, setFfmpegConfigList] = useState([]);

  useEffect(() => {
    fetch("/api/v1/get_ffmpeg_config_list")
      .then((response) => response.json())
      .then((data) => {
        setFfmpegConfigList(data.ffmpeg_config_list);
        if (data.ffmpeg_config_list.length > 0) {
          props.setSelectedFfmpegConfig(data.ffmpeg_config_list[0]);
        }
      })
      .catch((error) => {
        alert("get_ffmpeg_config_list error: " + error);
      });
  }, []);

  return (
    <div className="ffmpeg-config">
      <select
        onChange={(event) => {
          props.setSelectedFfmpegConfig(
            ffmpegConfigList[event.target.selectedIndex]
          );
        }}
      >
        {ffmpegConfigList.map((ffmpegConfig) => (
          <option key={ffmpegConfig.name}>{ffmpegConfig.name}</option>
        ))}
      </select>
      <span className="warp-word">{props.selectedFfmpegConfig.args}</span>
    </div>
  );
}

export default FfmpegConfig;
