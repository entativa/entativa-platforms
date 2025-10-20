package com.entativa.ui.feed

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.entativa.R

// MARK: - Post Menu Button (3-dot button on posts)
@Composable
fun PostMenuButton(
    post: Post,
    isOwnPost: Boolean,
    onEditPost: () -> Unit = {},
    onDeletePost: () -> Unit = {},
    onPinPost: () -> Unit = {},
    onArchivePost: () -> Unit = {},
    onToggleComments: () -> Unit = {},
    onReportPost: () -> Unit = {},
    onBlockUser: () -> Unit = {},
    onHidePost: () -> Unit = {},
    onNotInterested: () -> Unit = {},
    onShare: () -> Unit = {},
    onCopyLink: () -> Unit = {}
) {
    var showMenu by remember { mutableStateOf(false) }
    var showReportDialog by remember { mutableStateOf(false) }
    var showBlockDialog by remember { mutableStateOf(false) }
    var showDeleteDialog by remember { mutableStateOf(false) }
    
    // Menu Button
    IconButton(
        onClick = { showMenu = true },
        modifier = Modifier.size(32.dp)
    ) {
        Icon(
            painter = painterResource(R.drawable.ic_more_vertical),
            contentDescription = "Post options",
            tint = Color.Gray
        )
    }
    
    // Menu Dropdown
    DropdownMenu(
        expanded = showMenu,
        onDismissRequest = { showMenu = false }
    ) {
        if (isOwnPost) {
            // Own post options
            DropdownMenuItem(
                text = { Text("Edit Post") },
                onClick = {
                    onEditPost()
                    showMenu = false
                },
                leadingIcon = {
                    Icon(painterResource(R.drawable.ic_edit), null)
                }
            )
            
            DropdownMenuItem(
                text = { Text("Delete Post", color = Color.Red) },
                onClick = {
                    showDeleteDialog = true
                    showMenu = false
                },
                leadingIcon = {
                    Icon(painterResource(R.drawable.ic_trash), null, tint = Color.Red)
                }
            )
            
            Divider()
            
            DropdownMenuItem(
                text = { Text("Pin to Profile") },
                onClick = {
                    onPinPost()
                    showMenu = false
                },
                leadingIcon = {
                    Icon(painterResource(R.drawable.ic_pin), null)
                }
            )
            
            DropdownMenuItem(
                text = { Text("Archive Post") },
                onClick = {
                    onArchivePost()
                    showMenu = false
                },
                leadingIcon = {
                    Icon(painterResource(R.drawable.ic_archive), null)
                }
            )
            
            DropdownMenuItem(
                text = { Text(if (post.commentsDisabled) "Enable Comments" else "Turn Off Comments") },
                onClick = {
                    onToggleComments()
                    showMenu = false
                },
                leadingIcon = {
                    Icon(painterResource(R.drawable.ic_comment), null)
                }
            )
            
            Divider()
        } else {
            // Others' post options
            DropdownMenuItem(
                text = { Text("Report Post", color = Color.Red) },
                onClick = {
                    showReportDialog = true
                    showMenu = false
                },
                leadingIcon = {
                    Icon(painterResource(R.drawable.ic_flag), null, tint = Color.Red)
                }
            )
            
            DropdownMenuItem(
                text = { Text("Block @${post.author.username}", color = Color.Red) },
                onClick = {
                    showBlockDialog = true
                    showMenu = false
                },
                leadingIcon = {
                    Icon(painterResource(R.drawable.ic_block), null, tint = Color.Red)
                }
            )
            
            Divider()
            
            DropdownMenuItem(
                text = { Text("Hide Post") },
                onClick = {
                    onHidePost()
                    showMenu = false
                },
                leadingIcon = {
                    Icon(painterResource(R.drawable.ic_eye_off), null)
                }
            )
            
            DropdownMenuItem(
                text = { Text("Not Interested") },
                onClick = {
                    onNotInterested()
                    showMenu = false
                },
                leadingIcon = {
                    Icon(painterResource(R.drawable.ic_thumbs_down), null)
                }
            )
            
            Divider()
        }
        
        // Common options
        DropdownMenuItem(
            text = { Text("Share") },
            onClick = {
                onShare()
                showMenu = false
            },
            leadingIcon = {
                Icon(painterResource(R.drawable.ic_share), null)
            }
        )
        
        DropdownMenuItem(
            text = { Text("Copy Link") },
            onClick = {
                onCopyLink()
                showMenu = false
            },
            leadingIcon = {
                Icon(painterResource(R.drawable.ic_link), null)
            }
        )
    }
    
    // Report Dialog
    if (showReportDialog) {
        ReportPostDialog(
            post = post,
            onDismiss = { showReportDialog = false },
            onSubmit = { reason, details ->
                onReportPost()
                showReportDialog = false
            }
        )
    }
    
    // Block Confirmation Dialog
    if (showBlockDialog) {
        AlertDialog(
            onDismissRequest = { showBlockDialog = false },
            icon = {
                Icon(
                    painterResource(R.drawable.ic_block),
                    contentDescription = null,
                    tint = Color.Red,
                    modifier = Modifier.size(48.dp)
                )
            },
            title = {
                Text("Block @${post.author.username}?")
            },
            text = {
                Text("You won't see their posts and they won't be able to find your profile, posts or story on Entativa.")
            },
            confirmButton = {
                Button(
                    onClick = {
                        onBlockUser()
                        showBlockDialog = false
                    },
                    colors = ButtonDefaults.buttonColors(containerColor = Color.Red)
                ) {
                    Text("Block")
                }
            },
            dismissButton = {
                TextButton(onClick = { showBlockDialog = false }) {
                    Text("Cancel")
                }
            }
        )
    }
    
    // Delete Confirmation Dialog
    if (showDeleteDialog) {
        AlertDialog(
            onDismissRequest = { showDeleteDialog = false },
            icon = {
                Icon(
                    painterResource(R.drawable.ic_trash),
                    contentDescription = null,
                    tint = Color.Red,
                    modifier = Modifier.size(48.dp)
                )
            },
            title = {
                Text("Delete Post?")
            },
            text = {
                Text("This post will be permanently deleted. This action cannot be undone.")
            },
            confirmButton = {
                Button(
                    onClick = {
                        onDeletePost()
                        showDeleteDialog = false
                    },
                    colors = ButtonDefaults.buttonColors(containerColor = Color.Red)
                ) {
                    Text("Delete")
                }
            },
            dismissButton = {
                TextButton(onClick = { showDeleteDialog = false }) {
                    Text("Cancel")
                }
            }
        )
    }
}

// MARK: - Report Post Dialog
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ReportPostDialog(
    post: Post,
    onDismiss: () -> Unit,
    onSubmit: (ReportReason, String) -> Unit
) {
    var selectedReason by remember { mutableStateOf<ReportReason?>(null) }
    var additionalInfo by remember { mutableStateOf("") }
    
    AlertDialog(
        onDismissRequest = onDismiss
    ) {
        Surface(
            shape = MaterialTheme.shapes.large,
            tonalElevation = 6.dp
        ) {
            Column(
                modifier = Modifier.padding(24.dp)
            ) {
                // Title
                Text(
                    "Report Post",
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Bold
                )
                
                Spacer(modifier = Modifier.height(8.dp))
                
                Text(
                    "Help us understand what's happening",
                    fontSize = 14.sp,
                    color = Color.Gray
                )
                
                Spacer(modifier = Modifier.height(24.dp))
                
                // Report Reasons
                Column(
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    ReportReason.values().forEach { reason ->
                        Surface(
                            onClick = { selectedReason = reason },
                            modifier = Modifier.fillMaxWidth(),
                            shape = MaterialTheme.shapes.medium,
                            color = if (selectedReason == reason) Color(0xFF007CFC).copy(alpha = 0.1f) 
                                   else MaterialTheme.colorScheme.secondaryContainer
                        ) {
                            Row(
                                modifier = Modifier.padding(12.dp),
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                RadioButton(
                                    selected = selectedReason == reason,
                                    onClick = { selectedReason = reason }
                                )
                                
                                Spacer(modifier = Modifier.width(8.dp))
                                
                                Column {
                                    Text(
                                        reason.title,
                                        fontSize = 15.sp,
                                        fontWeight = FontWeight.Medium
                                    )
                                    Text(
                                        reason.subtitle,
                                        fontSize = 12.sp,
                                        color = Color.Gray
                                    )
                                }
                            }
                        }
                    }
                }
                
                // Additional Info
                if (selectedReason != null) {
                    Spacer(modifier = Modifier.height(16.dp))
                    
                    Text(
                        "Additional Information (Optional)",
                        fontSize = 13.sp,
                        color = Color.Gray
                    )
                    
                    Spacer(modifier = Modifier.height(8.dp))
                    
                    OutlinedTextField(
                        value = additionalInfo,
                        onValueChange = { additionalInfo = it },
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(100.dp),
                        placeholder = { Text("Describe the issue...") }
                    )
                }
                
                Spacer(modifier = Modifier.height(24.dp))
                
                // Buttons
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.End
                ) {
                    TextButton(onClick = onDismiss) {
                        Text("Cancel")
                    }
                    
                    Spacer(modifier = Modifier.width(8.dp))
                    
                    Button(
                        onClick = {
                            selectedReason?.let {
                                onSubmit(it, additionalInfo)
                            }
                        },
                        enabled = selectedReason != null,
                        colors = ButtonDefaults.buttonColors(containerColor = Color.Red)
                    ) {
                        Text("Submit Report")
                    }
                }
            }
        }
    }
}

// MARK: - Report Reasons
enum class ReportReason(val title: String, val subtitle: String) {
    SPAM("Spam", "Unwanted commercial content or spam"),
    INAPPROPRIATE("Inappropriate Content", "Nudity, sexual content, or graphic violence"),
    HARASSMENT("Harassment or Bullying", "Bullying or harassment of others"),
    HATE_SPEECH("Hate Speech", "Hateful or discriminatory content"),
    VIOLENCE("Violence", "Content promoting violence or terrorism"),
    FALSE_INFO("False Information", "Misleading or false information"),
    SCAM("Scam or Fraud", "Fraudulent content or scams"),
    OTHER("Something Else", "Something not listed here")
}

// MARK: - Post Extension
val Post.commentsDisabled: Boolean
    get() = false // TODO: real data from API
