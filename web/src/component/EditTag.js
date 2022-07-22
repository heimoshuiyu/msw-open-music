import { useState, useEffect, useContext } from "react";
import { useParams, useNavigate } from "react-router";
import { tr, Tr, langCodeContext } from "../translate";

function EditTag() {
  let params = useParams();
  let navigate = useNavigate();
  const { langCode } = useContext(langCodeContext);

  const [tag, setTag] = useState({
    id: "",
    name: "",
    description: "",
    created_by_user: {
      id: "",
      username: "",
      role: "",
      avatar_id: "",
    },
  });

  function refreshTagInfo() {
    fetch("/api/v1/get_tag_info", {
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
          setTag(data.tag);
        }
      });
  }

  function updateTagInfo() {
    fetch("/api/v1/update_tag", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: parseInt(params.id),
        name: tag.name,
        description: tag.description,
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          alert(tr("Tag updated successfully", langCode));
          refreshTagInfo();
        }
      });
  }

  useEffect(() => {
    refreshTagInfo();
  }, []);

  function deleteTag() {
    fetch("/api/v1/delete_tag", {
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
          alert(tr("Tag deleted successfully", langCode));
          navigate(-1);
        }
      });
  }

  return (
    <div className="page">
      <h3>{Tr("Edit Tag")}</h3>
      <div>
        <label htmlFor="id">{Tr("ID")}</label>
        <input
          type="text"
          disabled
          name="id"
          id="id"
          value={tag.id}
          onChange={(e) => setTag({ ...tag, id: e.target.value })}
        />
        <label htmlFor="name">{Tr("Created by")}</label>
        <input
          type="text"
          disabled
          name="created_by_user_username"
          id="created_by_user_username"
          value={tag.created_by_user.username}
          onChange={(e) =>
            setTag({
              ...tag,
              created_by_user: {
                ...tag.created_by_user,
                username: e.target.value,
              },
            })
          }
        />
        <label htmlFor="name">{Tr("Name")}</label>
        <input
          type="text"
          name="name"
          id="name"
          value={tag.name}
          onChange={(e) => setTag({ ...tag, name: e.target.value })}
        />
        <label htmlFor="description">{Tr("Description")}</label>
        <textarea
          name="description"
          id="description"
          value={tag.description}
          onChange={(e) => setTag({ ...tag, description: e.target.value })}
        />
        <button onClick={deleteTag}>{Tr("Delete")}</button>
        <button onClick={() => updateTagInfo()}>{Tr("Save")}</button>
      </div>
    </div>
  );
}

export default EditTag;
