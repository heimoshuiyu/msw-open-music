# DBMS Group Project Problem Description

- Group 1

The Internet infrastructure construction has made the network speed development faster. With the fast Internet, people are gradually migrating various data and services to the cloud. For example, NetEase Cloud Music, Spotify, and Apple Music, we call them streaming media platforms. The definition of streaming media platform is that users purchase the digital copyright of music and then play the music online on the platform.

Generally speaking, users cannot buy music that is not available on the platform. The user cannot download the digital file of the music (the user purchases the right to play instead of the right to copy). Users cannot upload their music to the platform.

However, in the era of digital copyright, there are still many advantages to getting original music files, such as no need to install a dedicated player; free copying to other devices (without violating copyright); no risk of music unavailable from the platform; no play records and privacy will be tracked by the platform.

Some people don't like streaming platforms. They like to collect music (download or buy CDs) and save it on their computers. But as more and more music is collected (over 70,000 songs and in total size of 800GB), it becomes very difficult to manage files. It is difficult for them to find where the songs they want to listen to are saved. Also, lossless music files are large and difficult to play online.

So we decided to develop a project based on database knowledge to help people who have collected a lot of music to enjoy their music simply.

The features of the project we designed are as follows:

- Open. Independent front-end (GUI) and back-end (server program), using API to communicate.
- Easy to use. Minimize dependencies, allowing users to configure quickly and simply.
- Lightweight. The program is small in size and quick to install.
- High performance. Only do what should be done, no features that will lead to poor performance.
- Cross-platform. The project can run on computers, mobile phones, Linux, Windows, macOS, and X86 and ARM processor architectures.
- Extensibility. Access to cloud OSS (Object Storage Service), reverse proxy, or other external software.

Our project has the following functions:

- Index file. Index local files into the database.
- Search. Search for music based on name/album/tag/comment, sorted by rating or other columns.
- Play. Play music online, play music randomly and play music at a low bit rate on a bad network.
- User management. Users can register and log in.
- Comment. Users can give a like or comment on the music.
- Management. The administrator can upload music, update or delete the database.
- Share. Generate a link to share the music with others.

After research and discussion, in order to meet the above requirements, we decided to use the Golang programming language on the backend. SQLite as a database program. Vue as the front-end GUI interface.