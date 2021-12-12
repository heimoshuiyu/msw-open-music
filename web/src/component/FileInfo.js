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
  const [tags, setTags] = useState([]);
  const [tagsOnFile, setTagsOnFile] = useState([]);
  const [selectedTagID, setSelectedTagID] = useState("");

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

  useEffect(() => {
    refresh();
    getTags();
    getTagsOnFile();
  }, []);

  return (
    <div className="page">
      <h3>File Details</h3>
      <div>
        <button>Download</button>
        <button
          onClick={() => {
            props.setPlayingFile(file);
          }}
        >
          Play
        </button>
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
      <div>
        <label>Tags:</label>
        <ul>
          {tagsOnFile.map((tag) => {
            return (
              <li key={tag.id}>
                <button>{tag.name}</button>
                <button>Remove</button>
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
            <option value="">Select a tag</option>
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
                alert("Please select a tag");
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
            Add Tag
          </button>
        </div>
      </div>
    </div>
  );
}

export default FileInfo;
