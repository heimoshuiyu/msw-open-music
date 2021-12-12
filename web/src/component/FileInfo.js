import { useNavigate, useParams } from "react-router";
import { useEffect, useState } from "react";

function FileInfo(props) {
  let navigate = useNavigate();
  let params = useParams();
  const [file, setFile] = useState({
    id: "",
    folder_id: "",
    foldername: "",
    filename: "",
    filesize: "",
  });

  function refresh() {
    fetch(`/api/v1/get_file_info`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: parseInt(params.id),
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setFile(data);
        }
      });
  }

  useEffect(() => {
    refresh();
  }, []);

  return (
    <div className="page">
      <h3>File Details</h3>
      <div>
        <button>Download</button>
        <button onClick={() => {
          props.setPlayingFile(file);
        }}>Play</button>
        <button
          onClick={() => {
            navigate(`/files/${params.id}/share`);
          }}
        >
          Share
        </button>
      </div>
      <div>
        <label htmlFor="foldername">Folder Name:</label>
        <input
          type="text"
          id="foldername"
          value={file.foldername}
          onClick={() => {
            navigate(`/folders/${file.folder_id}`);
          }}
          readOnly
        />
        <label htmlFor="filename">File Name:</label>
        <input type="text" id="filename" value={file.filename} readOnly />
        <label htmlFor="filesize">File Size:</label>
        <input type="text" id="filesize" value={file.filesize} readOnly />
      </div>
    </div>
  );
}

export default FileInfo;
