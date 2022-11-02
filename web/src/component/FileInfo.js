import {useNavigate, useParams} from "react-router";
import {useContext, useEffect, useState} from "react";
import {Tr, tr, langCodeContext} from "../translate";

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
  const [tags, setTags] = useState([]);
  const [tagsOnFile, setTagsOnFile] = useState([]);
  const [selectedTagID, setSelectedTagID] = useState("");
  const {langCode} = useContext(langCodeContext);
  const [ffprobeInfo, setFfprobeInfo] = useState("");

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

  function getTags() {
    fetch(`/api/v1/get_tags`)
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setTags(data.tags);
        }
      });
  }

  function getTagsOnFile() {
    fetch(`/api/v1/get_tags_on_file`, {
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
          setTagsOnFile(data.tags);
        }
      });
  }

  function removeTagOnFile(tag_id) {
    fetch(`/api/v1/delete_tag_on_file`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        file_id: parseInt(params.id),
        tag_id: tag_id,
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          getTagsOnFile();
        }
      });
  }

  function deleteFile() {
    // show Warning
    if (
      window.confirm(tr("Are you sure you want to delete this file?", langCode))
    ) {
      fetch(`/api/v1/delete_file`, {
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
            navigate(-1);
          }
        });
    }
  }

  function updateFilename() {
    fetch(`/api/v1/update_filename`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: parseInt(params.id),
        filename: file.filename,
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          alert(tr("Filename updated", langCode));
          refresh();
        }
      });
  }

  function resetFilename() {
    fetch(`/api/v1/reset_filename`, {
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
          refresh();
        }
      });
  }

  useEffect(() => {
    refresh();
    getTags();
    getTagsOnFile();
  }, []);

  const downloadURL = "/api/v1/get_file_direct?id=" + file.id;

  return (
    <div className="page">
      <h3>{Tr("File Details")}</h3>
      <div>
        <a href={downloadURL} download>
          <button>{Tr("Download")}</button>
        </a>
        <button
          onClick={() => {
            props.setPlayingFile(file);
          }}
        >
          {Tr("Play")}
        </button>
        <button
          onClick={() => {
            navigate(`/files/${params.id}/review`);
          }}
        >
          {Tr("Review")}
        </button>
        <button
          onClick={() => {
            navigate(`/files/${params.id}/share`);
          }}
        >
          {Tr("Share")}
        </button>
        <button
          onClick={() => {
            deleteFile();
          }}
        >
          {Tr("Delete")}
        </button>
      </div>
      <div>
        <label htmlFor="foldername">{Tr("Folder Name")}</label>
        <input
          type="text"
          id="foldername"
          value={file.foldername}
          onClick={() => {
            navigate(`/folders/${file.folder_id}`);
          }}
          readOnly
        />
        <label htmlFor="filename">{Tr("Filename")}</label>
        <input
          type="text"
          id="filename"
          value={file.filename}
          onChange={(event) => {
            setFile({
              ...file,
              filename: event.target.value,
            });
          }}
        />
        <label htmlFor="filesize">{Tr("File size")}</label>
        <input type="text" id="filesize" value={file.filesize} readOnly />
      </div>
      <div className="horizontal">
        <button onClick={updateFilename}>{Tr("Save")}</button>
        <button onClick={resetFilename}>{Tr("Reset")}</button>
      </div>
      <div>
        <label>{Tr("Tags")}</label>
        <ul>
          {tagsOnFile.map((tag) => {
            return (
              <li key={tag.id}>
                <button
                  onClick={() => {
                    navigate(`/manage/tags/${tag.id}`);
                  }}
                >
                  {tag.name}
                </button>
                <button
                  onClick={() => {
                    removeTagOnFile(tag.id);
                  }}
                >
                  {Tr("Remove")}
                </button>
              </li>
            );
          })}
        </ul>
        <div>
          <select
            onChange={(e) => {
              setSelectedTagID(e.target.value);
            }}
          >
            <option value="">{tr("Select a tag", langCode)}</option>
            {tags.map((tag) => {
              return (
                <option key={tag.id} value={tag.id}>
                  {tag.name}
                </option>
              );
            })}
          </select>
          <button
            onClick={() => {
              // check empty
              if (selectedTagID === "") {
                alert(tr("Please select a tag", langCode));
                return;
              }
              fetch(`/api/v1/put_tag_on_file`, {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify({
                  file_id: parseInt(params.id),
                  tag_id: parseInt(selectedTagID),
                }),
              })
                .then((res) => res.json())
                .then((data) => {
                  if (data.error) {
                    alert(data.error);
                  } else {
                    getTagsOnFile();
                  }
                });
            }}
          >
            {Tr("Add tag")}
          </button>
        </div>
      </div>

      <button onClick={async () => {
        const resp = await fetch(`/api/v1/get_file_ffprobe_info`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            id: parseInt(params.id),
          }),
        });
        const text = await resp.text();
        setFfprobeInfo(text);
      }}>FFprobe</button>

      {ffprobeInfo && <textarea
        style={{
          height: "30em",
        }}
      >{ffprobeInfo}</textarea>}

    </div>
  );
}

export default FileInfo;
