#import <Cocoa/Cocoa.h>
#include "dlg.h"

void* NSStr(void* buf, int len) {
	return (void*)[[NSString alloc] initWithBytes:buf length:len encoding:NSUTF8StringEncoding];
}

void NSRelease(void* obj) {
	[(NSObject*)obj release];
}

@interface AlertDlg : NSObject {
	AlertDlgParams* params;
	DlgResult result;
}
+ (AlertDlg*)init:(AlertDlgParams*)params;
- (DlgResult)run;
@end

DlgResult alertDlg(AlertDlgParams* params) {
	return [[AlertDlg init:params] run];
}

@implementation AlertDlg
+ (AlertDlg*)init:(AlertDlgParams*)params {
	AlertDlg* d = [AlertDlg alloc];
	d->params = params;
	return d;
}

- (DlgResult)run {
	if(![NSThread isMainThread]) {
		[self performSelectorOnMainThread:@selector(run) withObject:nil waitUntilDone:YES];
		return self->result;
	}
	NSAlert* alert = [[NSAlert alloc] init];
	if(self->params->title != nil) {
		[[alert window] setTitle:[[NSString alloc] initWithUTF8String:self->params->title]];
	}
	[alert setMessageText:[[NSString alloc] initWithUTF8String:self->params->msg]];
	switch (self->params->style) {
	case MSG_YESNO:
		[alert addButtonWithTitle:@"Yes"];
		[alert addButtonWithTitle:@"No"];
		break;
	case MSG_ERROR:
		[alert setIcon:[NSImage imageNamed:NSImageNameCaution]];
		[alert addButtonWithTitle:@"OK"];
		break;
	case MSG_INFO:
		[alert setIcon:[NSImage imageNamed:NSImageNameInfo]];
		[alert addButtonWithTitle:@"OK"];
		break;
	}
	self->result = [alert runModal] == NSAlertFirstButtonReturn ? DLG_OK : DLG_CANCEL;
	return self->result;
}
@end

@interface FileDlg : NSObject {
	FileDlgParams* params;
	DlgResult result;
}
+ (FileDlg*)init:(FileDlgParams*)params;
- (DlgResult)run;
@end

DlgResult fileDlg(FileDlgParams* params) {
	return [[FileDlg init:params] run];
}

@implementation FileDlg
+ (FileDlg*)init:(FileDlgParams*)params {
	FileDlg* d = [FileDlg alloc];
	d->params = params;
	return d;
}

- (DlgResult)run {
	if(![NSThread isMainThread]) {
		[self performSelectorOnMainThread:@selector(run) withObject:nil waitUntilDone:YES];
	} else if(self->params->mode == SAVEDLG) {
		self->result = [self save];
	} else {
		self->result = [self load];
	}
	return self->result;
}

- (NSInteger)runPanel:(NSSavePanel*)panel {
	[panel setFloatingPanel:YES];
	if(self->params->title != nil) {
		[panel setTitle:[[NSString alloc] initWithUTF8String:self->params->title]];
	}
	if(self->params->numext > 0) {
		[panel setAllowedFileTypes:[NSArray arrayWithObjects:(NSString**)self->params->exts count:self->params->numext]];
	}
	if(self->params->relaxext) {
		[panel setAllowsOtherFileTypes:YES];
	}
	if(self->params->startDir) {
		[panel setDirectoryURL:[NSURL URLWithString:[[NSString alloc] initWithUTF8String:self->params->startDir]]];
	}
	if(self->params->filename != nil) {
		[panel setNameFieldStringValue:[[NSString alloc] initWithUTF8String:self->params->filename]];
	}
	return [panel runModal];
}

- (DlgResult)save {
	NSSavePanel* panel = [NSSavePanel savePanel];
	if(![self runPanel:panel]) {
		return DLG_CANCEL;
	} else if(![[panel URL] getFileSystemRepresentation:self->params->buf maxLength:self->params->single_buf]) {
		return DLG_URLFAIL;
	}
	return DLG_OK;
}

- (DlgResult)load {
	NSOpenPanel* panel = [NSOpenPanel openPanel];
	if(self->params->mode == DIRDLG) {
		[panel setCanChooseDirectories:YES];
		[panel setCanChooseFiles:NO];
	}

	if (self->params->multiple == 1) {
		[panel setAllowsMultipleSelection:YES];
	}

	if(![self runPanel:panel]) {
		return DLG_CANCEL;
	}

	NSArray* urls = [panel URLs];

	if ([urls count] > self->params->max_files) {
		return DLG_TOOMANY;
	}

	char* bufp = self->params->buf;
	size_t written = 0;

	for (int i = 0; i < [urls count]; i++) {
		NSURL* url = [urls objectAtIndex:i];
		if(![url getFileSystemRepresentation:bufp maxLength:self->params->single_buf]) {
			return DLG_URLFAIL;
		}

		size_t newLen = strlen(bufp);
		bufp += (newLen + 1);
		written += (newLen + 1);
	}

	if (written > 0) {
		written -= 1;
	}

	*self->params->written = written;
	return DLG_OK;
}

@end