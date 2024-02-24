package projects

import "mime/multipart"

func validateAttachment(marketingFiles, routeFiles, eventMapFiles, eventDetailsFiles, screenshotFiles []*multipart.FileHeader) error {
	if len(marketingFiles) == 0 {
		return &MarketingFilesRequiredError{}
	}
	if len(routeFiles) == 0 {
		return &RouteFilesRequiredError{}
	}
	if len(eventMapFiles) == 0 {
		return &EventMapFilesRequiredError{}
	}
	if len(eventDetailsFiles) == 0 {
		return &EventDetailsFilesRequiredError{}
	}
	if len(screenshotFiles) == 0 {
		return &ScreenshotFilesRequiredError{}
	}
	return nil
}
