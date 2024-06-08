package user

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/go-oss/image/imageutil"
	"github.com/nfnt/resize"
)

/*
* @Internal
* Update user with new name & private profile flag.
 */
func updateUserNameAndPrivateProfile(
	user *models.User,
	name string,
	privateProfile bool,
) error {
	// 1. Check whether new user name is valid.
	if !isValidUserName(name) && user.Name != name {
		return utils.MakeErrorWithCode(
			"user_update",
			"updateUserNameAndPrivateProfile",
			"invalid user name",
			ErrCodeInvalidUserName,
			fmt.Errorf(
				"oldName: %s, newName: %s",
				user.Name,
				name,
			),
		)
	}

	// 2. Save user info.
	if err := saveUser(
		user,
		map[string]interface{}{
			"name":            name,
			"private_profile": privateProfile,
		},
	); err != nil {
		return utils.MakeError(
			"user_update",
			"updateUserNameAndPrivateProfile",
			"failed to save update user info",
			err,
		)
	}
	return nil
}

/*
* @Internal
* Update user avatar.
 */
func updateUserAvatar(
	user *models.User,
	imageHeader *multipart.FileHeader,
) (string, error) {
	// 1. Read image data from multipart file header.
	fileType, image, err := readImageData(imageHeader)
	if err != nil {
		return "", utils.MakeError(
			"user_update",
			"updateUserAvatar",
			"failed to read image data",
			err,
		)
	}

	// 2. Get image key from image file.
	key := getKeyFromImage(*image)

	// 3. Delete old image if is not default avatar.
	if !isDefaultAvatar(user.Avatar) {
		if err := utils.DeleteImage(
			getKeyfromUrl(user.Avatar),
		); err != nil {
			return "", utils.MakeError(
				"user_update",
				"updateUserAvatar",
				"failed to delete old image",
				err,
			)
		}
	}

	// 4. Upload new user avatar.
	url, err := utils.UploadImage(*image, fileType, key)
	if err != nil {
		return "", utils.MakeError(
			"user_update",
			"updateUserAvatar",
			"failed to upload new avatar",
			err,
		)
	}

	// 5. Save user info with new avatar.
	if err := saveUser(
		user,
		map[string]interface{}{
			"avatar": url,
		},
	); err != nil {
		return "", utils.MakeError(
			"user_update",
			"updateUserAvatar",
			"failed to save user info",
			err,
		)
	}
	return url, nil
}

/*
* @Internal
* Read image buffer from multipart file reader.
 */
func readImageData(
	imageHeader *multipart.FileHeader,
) (string, *image.Image, error) {
	// 1. Get multipart file from image header.
	multipartFile, err := imageHeader.Open()
	if err != nil {
		return "", nil, utils.MakeError(
			"user_update",
			"readImageData",
			"failed to open multipart file",
			err,
		)
	}
	defer multipartFile.Close()

	// 2. Get image size.
	size := imageHeader.Size

	// 3. Read image to buffer.
	buffer := make([]byte, size)
	multipartFile.Read(buffer)

	// 4. Get image file type.
	fileType := http.DetectContentType(buffer)

	// 5. Move seek to the head of file
	if newSeek, err := multipartFile.Seek(0, 0); err != nil {
		return "", nil, utils.MakeError(
			"user_update",
			"readImageData",
			"failed to move file seek to head",
			err,
		)
	} else {
		fmt.Println("New Seek:", newSeek)
	}

	fmt.Println("File Type:", fileType)
	fmt.Println("File Size:", size)
	// 6. Decode image.
	originalImage, err := decodeImage(
		fileType,
		multipartFile,
	)
	if err != nil {
		return "", nil, utils.MakeError(
			"user_update",
			"readImageData",
			"failed to decode original image",
			err,
		)
	}

	// 7. Resize original image and get new image.
	new_image := resize.Resize(
		100,
		100,
		originalImage,
		resize.Lanczos2,
	)
	return fileType, &new_image, nil
}

/*
* @Internal
* Decode and return original image.
 */
func decodeImage(fileType string, multipartFile multipart.File) (image.Image, error) {
	originalImage, err := imageutil.Decode(
		multipartFile,
	)
	if err != nil {
		return nil, utils.MakeError(
			"user_update",
			"readImageData",
			"failed to decode original image",
			err,
		)
	} else {
		return originalImage.Image, nil
	}
}

/*
* @Internal
* Get key from image object.
 */
func getKeyFromImage(img image.Image) string {
	str := fmt.Sprintf("%v", img)
	sum := sha256.Sum256([]byte(str))
	return "avatar/" + hex.EncodeToString(sum[:]) + ".png"
}

/*
* @Internal
* Get key from image url.
 */
func getKeyfromUrl(url string) string {
	parts := strings.Split(url, "/")
	return "avatar/" + parts[len(parts)-1]
}

/*
* @Internal
* Check whether default avatar url.
 */
func isDefaultAvatar(imgUrl string) bool {
	return imgUrl == DEFAULT_USER_AVATAR_URL
}
